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

const (
	GetUserInfo        = "GET USER INFO"
	RegisterUser       = "REGISTER USER"
	UpdateUserLocation = "UPDATE USER LOCATION"
	DeleteUser         = "DELETE USER"
	RequestAssociation = "REQUEST ASSOCIATION"
	RequestDissociation = "REQUEST DISSOCIATION"
	RequestLockChange  = "REQUEST LOCK CHANGE"

	GetSpidInfo        = "GET SPID INFO"
	RegisterSpid       = "REGISTER SPID"
	UpdateBatteryInfo  = "UPDATE BATTERY INFO"
	UpdateSpidLocation = "UPDATE SPID LOCATION"
	DeleteSpid         = "DELETE SPID"
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

func NewHandler(basePath string) Handler {
	return Handler{db.NewManager(basePath)}
}

func DefaultResponse(request Request) Response {
	return Response{
		ID:   request.ID,
		Type: "<RESPONSE>:" + request.Type,
		Ok:   false,
		Body: map[string]interface{}{},
	}
}

func (h Handler) requestProcessor(incomingRequests chan string, outgoingResponses chan string) {
	for {
		incomingRequest := <-incomingRequests
		response, _ := h.ProcessRequest(incomingRequest)
		outgoingResponses <- response
	}
}

func (h Handler) Listen(port string) {
	incomingRequests := make(chan string)
	outgoingResponses := make(chan string)
	go h.requestProcessor(incomingRequests, outgoingResponses)
	tcpServer.Listen(port, incomingRequests, outgoingResponses)
}

func invalidRequest(request Request) (response Response, ok bool) {
	response = DefaultResponse(request)
	response.Body["message"] = fmt.Sprintf("Invalid request type '%s'", request.Type)
	return response, false
}

func checkKeys(m map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if m[key] == nil {
			return key
		}
	}
	return ""
}

func (h Handler) ProcessRequest(incomingRequest string) (jsonResponse string, ok bool) {
	var request Request
	err := json.Unmarshal([]byte(incomingRequest), &request)

	if err != nil {
		errorMessage := fmt.Sprintf("Invalid request: `%s`\n`%s`", incomingRequest, err)
		log.Print(errorMessage)
		marshaledResponse, err := json.Marshal(Response{
			ID:   uuid.Nil,
			Type: "<RESPONSE>:INVALID",
			Ok:   false,
			Body: map[string]interface{}{"message": errorMessage},
		})
		eh.HandleFatal(err)
		return string(marshaledResponse), false
	}

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
	case RequestAssociation:
		handler = h.requestAssociation
	case RequestDissociation:
		handler = h.requestDissociation
	case RequestLockChange:
		handler = h.requestLockChange
	case GetSpidInfo:
		handler = h.getSpidInfo
	case RegisterSpid:
		handler = h.registerSpid
	case UpdateBatteryInfo:
		handler = h.updateBatteryInfo
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
	log.Printf("Response: `%s`", response)
	marshaledResponse, err := json.Marshal(response)
	eh.HandleFatal(err)
	return string(marshaledResponse), ok
}
