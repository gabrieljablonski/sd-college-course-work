package requestHandling

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ActiveState/tail"
	"log"
	"main/db"
	eh "main/errorHandling"
	"os"
	"time"
)

const (
	GetUserInfo        = "GET USER INFO"
	RegisterUser       = "REGISTER USER"
	UpdateUserLocation = "UPDATE USER LOCATION"
	DeleteUser         = "DELETE USER"
	RequestAssociation = "REQUEST ASSOCIATION"
	RequestDissociation = "REQUEST DISSOCIATION"
	RequestSpidInfo    = "REQUEST SPID INFO"
	RequestLockChange  = "REQUEST LOCK CHANGE"

	GetSpidInfo        = "GET SPID INFO"
	RegisterSpid       = "REGISTER SPID"
	UpdateBatteryInfo  = "UPDATE BATTERY INFO"
	UpdateSpidLocation = "UPDATE SPID LOCATION"
	DeleteSpid         = "DELETE SPID"

	TimedOut = "TIMEOUT"
	Invalid = "INVALID"

	DefaultLogPath     = "requestHandling/request_logs.spd"
	DefaultMaxBufferedRequests = 100
	DefaultWriteToFilePeriod   = 5000*time.Second
)

type GenericMessage struct {
	Message string `json:"message"`
	Received time.Time `json:"received"`
	Sum [16]byte   `json:"sum"`
}

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
	Manager           db.Manager
	OutGoingResponses chan GenericMessage
	Logger            *log.Logger
}

func NewHandler(basePath string) Handler {
	//TODO: if file is not empty, process requests
	logFile, err := os.OpenFile(DefaultLogPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	eh.HandleFatal(err)
	h := Handler{
		Manager: db.NewManager(basePath),
		OutGoingResponses: make(chan GenericMessage, DefaultMaxBufferedRequests),
		Logger: log.New(logFile, "", 0),
	}
	h.Manager.LoadFromFile()
	go h.WriteToFilePeriodically(DefaultWriteToFilePeriod)
	go h.handleLoggedRequests()
	return h
}

func (h Handler) WriteToFilePeriodically(period time.Duration) {
	for {
		time.Sleep(period)
		//TODO: determine which was the last processed request
		log.Print("Writing users and spids to file.")
		h.Manager.WriteUsersToFile()
		h.Manager.WriteSpidsToFile()
	}
}

func defaultResponse(request Request) Response {
	return Response{
		ID:   request.ID,
		Type: "<RESPONSE>:" + request.Type,
		Ok:   false,
		Body: map[string]interface{}{},
	}
}

func invalidResponse(request Request) (response Response, ok bool) {
	request.Type = Invalid
	response = defaultResponse(request)
	response.Body["message"] = fmt.Sprintf("Invalid request type '%s'", request.Type)
	return response, false
}

func (h Handler) QueueRequest(requestMessage GenericMessage) error {
	marshaledRequest, err := json.Marshal(requestMessage)
	eh.HandleFatal(err)
	err = h.Logger.Output(2, string(marshaledRequest))
	f, ferr := os.Open(DefaultLogPath)
	eh.HandleFatal(ferr)
	ferr = f.Close()
	eh.HandleFatal(ferr)
	return err
}

func (h *Handler) cleanUp(discardedRequestSum [16]byte) {
	for {
		responseMessage := <-h.OutGoingResponses
		if responseMessage.Sum == discardedRequestSum {
			return
		}
		h.OutGoingResponses <- responseMessage
	}
}

func (h *Handler) GetResponse(requestMessage GenericMessage, timeout time.Duration) (GenericMessage, error) {
	for {
		select {
		case responseMessage := <-h.OutGoingResponses:
			if responseMessage.Sum == requestMessage.Sum {
				return responseMessage, nil
			}
			h.OutGoingResponses <- responseMessage
		case <-time.After(timeout):
			go h.cleanUp(requestMessage.Sum)
			marshaledResponse, err := json.Marshal(Response{
				ID:   uuid.Nil,
				Type: "<RESPONSE>:" + TimedOut,
				Ok:   false,
				Body: map[string]interface{}{"message": "Request timed out."},
			})
			eh.HandleFatal(err)
			return GenericMessage{
				Message: string(marshaledResponse),
				Sum:     requestMessage.Sum,
			}, fmt.Errorf("request timed out")
		}
	}
}

func (h *Handler) handleLoggedRequests() {
	logPath := h.Manager.FileManager.BasePath + string(os.PathSeparator) + DefaultLogPath
	logTail, err := tail.TailFile(logPath, tail.Config{
		Location:    nil,
		ReOpen:      false,
		MustExist:   false,
		Poll:        true,
		Pipe:        false,
		RateLimiter: nil,
		Follow:      true,
		MaxLineSize: 0,
		Logger:      nil,
	})
	eh.HandleFatal(err)
	for {
		for line := range logTail.Lines {
			log.Printf("Processing logged request: `%s`", line.Text)
			var requestMessage GenericMessage
			err := json.Unmarshal([]byte(line.Text), &requestMessage)
			eh.HandleFatal(err)
			jsonResponse, _ := h.processRequest(requestMessage.Message)

			h.OutGoingResponses <- GenericMessage{
				Message:  jsonResponse,
				Received: requestMessage.Received,
				Sum:      requestMessage.Sum,
			}
		}
	}
}

func (h Handler) processRequest(incomingRequest string) (jsonResponse string, ok bool) {
	var request Request
	err := json.Unmarshal([]byte(incomingRequest), &request)

	if err != nil {
		errorMessage := fmt.Sprintf("Invalid request: `%s`\n`%s`", incomingRequest, err)
		log.Print(errorMessage)
		marshaledResponse, err := json.Marshal(Response{
			ID:   uuid.Nil,
			Type: "<RESPONSE>:" + Invalid,
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
	case RequestSpidInfo:
		handler = h.requestSpidInfo
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
		handler = invalidResponse
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
