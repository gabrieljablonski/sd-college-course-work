package requestHandling

import (
	"fmt"
	"github.com/google/uuid"
	"main/entities"
	"main/gps"
	"main/utils"
	"time"
)

func (h *Handler) querySpid(spidID string) (entities.Spid, error) {
	id, err := uuid.Parse(spidID)
	if err != nil {
		return entities.Spid{}, fmt.Errorf("invalid spid id: %s", err)
	}
	spid, err := h.Manager.QuerySpid(id)
	if err != nil {
		return entities.Spid{}, err
	}
	return spid, nil
}

func (h *Handler) getSpidInfo(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	if request.Body["spid_id"] == nil {
		response.Body["message"] = "Missing spid id."
		return response, false
	}
	spid, err := h.querySpid(request.Body["spid_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	response.Ok = true
	response.Body["spid"] = spid
	return response, true
}

func (h *Handler) registerSpid(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	spid := entities.NewSpid()
	err := h.Manager.RegisterSpid(spid)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to register spid: %s", err)
		return response, false
	}
	response.Body["spid"] = spid
	response.Ok = true
	return response, true
}

func (h *Handler) updateBatteryInfo(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	missingKey := utils.CheckKeys(request.Body, []string{"spid_id", "battery_level"})
	if missingKey != "" {
		response.Body["message"] = fmt.Sprintf("Missing key: `%s`.", missingKey)
		return response, false
	}
	spid, err := h.querySpid(request.Body["spid_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	batteryLevel, okBattery := request.Body["battery_level"].(int)
	if !okBattery || (batteryLevel < 0 || batteryLevel > 100) {
		response.Body["message"] = "Invalid battery value."
		return response, false
	}
	spid.BatteryLevel = uint8(batteryLevel)
	err = h.Manager.UpdateSpid(spid)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update spid: `%s`", err)
		return response, false
	}
	response.Body["message"] = "Battery level updated."
	response.Body["spid"] = spid
	response.Ok = true
	return response, true
}

func (h *Handler) updateSpidLocation(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	missingKey := utils.CheckKeys(request.Body, []string{"location", "spid_id"})
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
	spid, err := h.querySpid(request.Body["spid_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	spid.Location = gps.GlobalPosition{
		Latitude:  latitude,
		Longitude: longitude,
	}
	spid.LastUpdated = time.Now()
	err = h.Manager.UpdateSpid(spid)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to update spid location: %s", err)
		return response, false
	}
	response.Body["message"] = "Location updated."
	response.Body["spid"] = spid
	response.Ok = true
	return response, true
}

func (h *Handler) deleteSpid(request Request) (response Response, ok bool) {
	response = defaultResponse(request)
	if request.Body["spid_id"] == nil {
		response.Body["message"] = "Missing spid id."
		return response, false
	}
	spid, err := h.querySpid(request.Body["spid_id"].(string))
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to validate spid id: %s", err)
		return response, false
	}
	err = h.Manager.DeleteSpid(spid)
	if err != nil {
		response.Body["message"] = fmt.Sprintf("Failed to delete spid: %s", err)
		return response, false
	}
	response.Body["message"] = fmt.Sprintf("Spid deleted.")
	response.Body["spid"] = spid
	response.Ok = true
	return response, true
}
