package requestHandling

import (
	"fmt"
	"github.com/google/uuid"
	"main/entities"
	"main/gps"
	"main/utils"
	"time"
)

func (h *Handler) queryUser(userID string) (entities.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return entities.User{}, fmt.Errorf("invalid user id: %s", err)
	}
	user, err := h.Manager.QueryUser(id)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (h *Handler) getUserInfo(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	if request.Body["user_id"] == nil {
		response.Body["message"] = "Missing user id."
		return response, false
	}
	user, err := h.queryUser(request.Body["user_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	response.Ok = true
	response.Body["user"] = user
	return response, true
}

func (h *Handler) registerUser(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	if request.Body["user_name"] == nil {
		response.Body["message"] = "Missing key `user_name`."
		return response, false
	}
	user := entities.NewUser(request.Body["user_name"].(string))
	err := h.Manager.RegisterUser(user)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to register user: %s", err)
		return response, false
	}
	response.Body["user"] = user
	response.Ok = true
	return response, true
}

func (h *Handler) updateUserLocation(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	missingKey := utils.CheckKeys(request.Body, []string{"location", "user_id"})
	if missingKey != "" {
		response.Body["message"] = fmt.Sprintf("Missing key: `%s`.", missingKey)
		return response, false
	}
	latitudeInterface := request.Body["location"].(map[string]interface{})["latitude"]
	longitudeInterface := request.Body["location"].(map[string]interface{})["longitude"]
	if latitudeInterface == nil || longitudeInterface == nil {
		response.Body["message"] = "Missing latitude or longitude."
		return response, false
	}
	latitude, okLat := latitudeInterface.(float64)
	longitude, okLong := longitudeInterface.(float64)
	if okLat && (latitude < -90 || latitude > 90) {
		okLat = false
	}
	if okLong && (longitude < -180 || longitude > 180) {
		okLong = false
	}
	if !(okLat && okLong) {
		response.Body["message"] = "Invalid latitude or longitude values."
		return response, false
	}
	user, err := h.queryUser(request.Body["user_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	user.Location = gps.GlobalPosition{
		Latitude:  latitude,
		Longitude: longitude,
	}
	user.LastUpdated = time.Now()
	err = h.Manager.UpdateUser(user)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update user location: %s", err)
		return response, false
	}
	response.Body["message"] = "Location updated."
	response.Body["user"] = user
	response.Ok = true
	return response, true
}

func (h Handler) deleteUser(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	if request.Body["user_id"] == nil {
		response.Body["message"] = "Missing user id."
		return response, false
	}
	user, err := h.queryUser(request.Body["user_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	err = h.Manager.DeleteUser(user)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to delete user: %s", err)
		return response, false
	}
	response.Body["message"] = fmt.Sprintf("User deleted.")
	response.Body["user"] = user
	response.Ok = true
	return response, true
}

//TODO: send message to spid with association/lock change

func (h *Handler) requestAssociation(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	missingKey := utils.CheckKeys(request.Body, []string{"user_id", "spid_id"})
	if missingKey != "" {
		response.Body["message"] = fmt.Sprintf("Missing key: `%s`.", missingKey)
		return response, false
	}
	userID := request.Body["user_id"].(string)
	spidID := request.Body["spid_id"].(string)
	user, err := h.queryUser(userID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	if user.CurrentSpidID != uuid.Nil {
		response.Body["message"] = fmt.Sprintf("User already associated to spid %s.", user.CurrentSpidID)
	}
	spid, err := h.querySpid(spidID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	if spid.CurrentUserID != uuid.Nil {
		response.Body["message"] = fmt.Sprintf("Spid already associated to user %s.", spid.CurrentUserID)
		return response, false
	}
	spid.CurrentUserID = user.ID
	user.CurrentSpidID = spid.ID
	err = h.Manager.UpdateUser(user)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update user info: `%s`", err)
		return response, false
	}
	err = h.Manager.UpdateSpid(spid)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update spid info: `%s`", err)
		return response, false
	}
	response.Body["message"] = fmt.Sprintf("User %s associated to spid %s.", user.ID, spid.ID)
	response.Body["user"] = user
	response.Ok = true
	return response, true
}

func (h *Handler) requestDissociation(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	if request.Body["user_id"] == nil {
		response.Body["message"] = "Missing user id."
		return response ,false
	}
	userID := request.Body["user_id"].(string)
	user, err := h.queryUser(userID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	if user.CurrentSpidID == uuid.Nil {
		response.Body["message"] = "User not currently associated to any spids."
		return response, false
	}
	spidID := user.CurrentSpidID
	user.CurrentSpidID = uuid.Nil
	err = h.Manager.UpdateUser(user)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to dissociate user: %s", err)
		return response, false
	}
	response.Body["message"] = fmt.Sprintf("User dissociated from spid %s.", spidID)
	response.Body["user"] = user
	response.Ok = true
	return response, true
}

func (h *Handler) requestSpidInfo(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	missingKey := utils.CheckKeys(request.Body, []string{"user_id", "spid_id"})
	if missingKey != "" {
		response.Body["message"] = fmt.Sprintf("Missing key: `%s`.", missingKey)
		return response, false
	}
	userID := request.Body["user_id"].(string)
	spidID := request.Body["spid_id"].(string)
	user, err := h.queryUser(userID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	spid, err := h.querySpid(spidID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	if user.CurrentSpidID != spid.ID || spid.CurrentUserID != user.ID {
		response.Body["message"] = fmt.Sprintf(
			"User with id %s not associated to spid with id %s.", user.ID, spid.ID)
		return response, false
	}
	response.Body["spid"] = map[string]interface{}{
		"id": spid.ID,
		"location": spid.Location,
		"battery_level": spid.BatteryLevel,
		"lock_state": spid.Lock}
	response.Ok = true
	return response, true
}

func (h *Handler) requestLockChange(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	missingKey := utils.CheckKeys(request.Body, []string{"user_id", "spid_id", "lock_state"})
	if missingKey != "" {
		response.Body["message"] = fmt.Sprintf("Missing key: `%s`.", missingKey)
		return response, false
	}
	lockState := request.Body["lock_state"].(string)
	userID := request.Body["user_id"].(string)
	spidID := request.Body["spid_id"].(string)
	if !entities.IsValidLockState(lockState) {
		response.Body["message"] = fmt.Sprintf("Invalid lock state: `%s`", lockState)
		return response, false
	}
	user, err := h.queryUser(userID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	spid, err := h.querySpid(spidID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	if user.CurrentSpidID != spid.ID || spid.CurrentUserID != user.ID {
		response.Body["message"] = fmt.Sprintf(
			"User with id %s not associated to spid with id %s.", user.ID, spid.ID)
		return response, false
	}
	err = spid.UpdateLockState(lockState, user.ID)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update spid lock state: %s", err)
		return response, false
	}
	err = h.Manager.UpdateSpid(spid)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update spid: %s", err)
		return response, false
	}
	// TODO: call back when spid checks the pending update
	response.Body["message"] = fmt.Sprintf("Lock state for spid %s updated to `%s`.", spid.ID, lockState)
	response.Body["spid"] = map[string]interface{}{
		"id": spid.ID,
		"location": spid.Location,
		"battery_level": spid.BatteryLevel,
		"lock_state": spid.Lock}
	response.Ok = true
	return response, true
}
