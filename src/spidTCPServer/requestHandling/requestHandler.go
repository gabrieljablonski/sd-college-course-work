package requestHandling

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"main/db"
	eh "main/errorHandling"
	"main/tcpServer"
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
	Body map[string]interface{} `json:"body"`
}

type Handler struct {
	Manager db.Manager
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

func NewHandler(basePath string) Handler {
	return Handler{db.NewManager(basePath)}
}

func DefaultResponse(request Request) Response {
	return Response{
		ID:   request.ID,
		Type: request.Type,
		Ok:   false,
		Body: map[string]interface{}{},
	}
}

func (h Handler) Listen(port string) {
	tcpServer.Listen(port)
}

func invalidRequest(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	response.Body["message"] = fmt.Sprintf("Invalid request type '%s'", request.Type)
	return response, false
}

func (h Handler) ProcessRequest(message string) (jsonResponse string, ok bool) {
	var request Request
	err := json.Unmarshal([]byte(message), &request)
	eh.HandleFatal(err)

	log.Printf("Processing request: %s", request)
	var handler func(Request) (Response, bool)

	switch request.Type {
	case GetUserInfo:
		handler = h.getUserInfo
	case RegisterUser:
		handler = h.registerUser
	case UpdateUserLocation:
		handler = h.updateUserLocation
	case DeleteUser:
		handler = h.deleteUser
	case RequestLockChange:
		handler = h.requestLockChange
	case GetSpidInfo:
		handler = h.getSpidInfo
	case RegisterSpid:
		handler = h.registerSpid
	case UpdateSpidLocation:
		handler = h.updateSpidLocation
	case DeleteSpid:
		handler = h.deleteSpid
	default:
		handler = invalidRequest
	}
	response, ok := handler(request)
	if !ok {
		log.Print("Request failed.\n")
	} else {
		log.Print("Request successful.\n")
	}
	log.Printf("Response: %s", response)
	marshaledResponse, err := json.Marshal(response)
	eh.HandleFatal(err)
	return string(marshaledResponse), ok
}
