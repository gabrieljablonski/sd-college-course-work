package requestHandling

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"main/entities"
	"main/gps"
	"time"
)

func (h Handler) queryUser(userID string) (entities.User, error) {
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

func (h Handler) getUserInfo(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	if request.Body["user_id"] == nil {
		response.Body["message"] = "Missing user id."
		return response, false
	}
	user, err := h.queryUser(request.Body["user_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate user id: %s", err)
		return response, false
	}
	log.Printf("Responding with user: %s", user)
	response.Ok = true
	response.Body["user"] = user
	return response, true
}

func (h Handler) registerUser(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	if request.Body["user_name"] == nil {
		response.Body["message"] = "Missing new user's name."
		return response, false
	}
	user := entities.NewUser(request.Body["user_name"].(string))
	err := h.Manager.RegisterUser(user)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to register user: %s", err)
		ok = false
	} else {
		response.Body["user"] = user
		response.Ok = true
		ok = true
	}
	return response, ok
}

func (h Handler) updateUserLocation(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	missingKey := checkKeys(request.Body, []string{"location", "user_id"})
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
		response.Body["message"] = fmt.Sprintf("Failed to update user: %s", err)
		return response, false
	}
	response.Body["message"] = "Location updated."
	response.Body["location"] = user.Location
	response.Ok = true
	return response, true
}

func (h Handler) deleteUser(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
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

func (h Handler) requestLockChange(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	missingKey := checkKeys(request.Body, []string{"user_id", "spid_id", "lock_state"})
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
	// TODO: call back when spid checks the pending update
	response.Body["message"] = fmt.Sprintf("Lock state for spid %s updated to `%s`.", spid.ID, lockState)
	response.Ok = true
	return response, false
}
