package requestHandling

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	eh "main/errorHandling"
)

type Request struct {
	ID   uuid.UUID              `json:"id"`
	Type string                 `json:"type"`
	Body map[string]interface{} `json:"body"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
	Type string `json:"type"`
	Ok bool `json:"ok"`
}

const (
	GetUserInfo        = "GET USER INFO"
	RegisterUser       = "REGISTER USER"
	UpdateUserLocation = "UPDATE USER LOCATION"
	DeleteUser         = "DELETE USER"
	RequestLockChange  = "REQUEST LOCK CHANGE"

	GetSpidInfo        = "GET SPID INFO"
	RegisterSpid       = "REGISTER SPID"
	UpdateSpidLocation = "UPDATE SPID LOCATION"
	DeleteSpid         = "DELETE SPID"
)

func ProcessRequest(message string) Response {
	var request Request
	err := json.Unmarshal([]byte(message), &request)
	eh.HandleFatal(err)

	var handler func(request Request) Response

	switch request.Type {
	case GetUserInfo:
		handler = getUserInfo
	case RegisterUser:
		handler = registerUser
	case UpdateUserLocation:
		handler = updateUserLocation
	case DeleteUser:
		handler = deleteUser
	case RequestLockChange:
		handler = requestLockChange
	case GetSpidInfo:
		handler = getSpidInfo
	case RegisterSpid:
		handler = registerSpid
	case UpdateSpidLocation:
		handler = updateSpidLocation
	case DeleteSpid:
		handler = deleteSpid
	default:
		log.Fatalf("Invalid request type '%s'", request.Type)
	}
	return handler(request)
}
