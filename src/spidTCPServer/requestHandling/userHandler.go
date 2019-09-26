package requestHandling

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"main/entities"
	eh "main/errorHandling"
)

func (h Handler) getUserInfo(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	if request.Body["id"] == nil {
		response.Body["message"] = fmt.Sprintf("Request missing user ID: %s", request)
		return response, false
	}
	id, err := uuid.Parse(request.Body["id"].(string))
	eh.HandleFatal(err)

	user := h.Manager.QueryUser(id)
	if user == (entities.User{}) {
		response.Body["message"] = "Request failed: user not found."
		return response, false
	}

	log.Printf("Responding with user: %s", user)
	response.Ok = true
	response.Body["user"] = user
	return response, true
}

func (h Handler) registerUser(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	if request.Body["name"] == nil {
		response.Body["message"] = "Missing new user's name."
		return response, false
	}
	user := entities.NewUser(request.Body["name"].(string))
	if !h.Manager.RegisterUser(user) {
		response.Body["message"] = "Failed to register user."
		ok = false
	} else {
		response.Body["user"] = user
		response.Ok = true
		ok = true
	}
	return response, ok
}

func (h Handler) updateUserLocation(request Request) (response Response, ok bool) {
	return Response{}, false
}

func (h Handler) deleteUser(request Request) (response Response, ok bool) {
	return Response{}, false
}

func (h Handler) requestLockChange(request Request) (response Response, ok bool) {
	return Response{}, false
}
