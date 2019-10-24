package requestHandling

import (
	"log"
	"os"
	"spidServer/requestHandling/grpc"
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

	DefaultLogPath     = "requestHandling" + string(os.PathSeparator) + "request_logs.spd"
	DefaultDirtyRequestsPath = "requestHandling" + string(os.PathSeparator) + "dirty_requests.spd"
	DefaultMaxBufferedRequests = 100
	DefaultWriteToFilePeriod   = 5000*time.Millisecond
)

type Handler struct {
	//Manager           db.Manager
	OutGoingResponses chan GenericMessage
	LoggerPending     *log.Logger
	LoggerDirty       *log.Logger
	WritingToFile     bool
	WritingToMemory   bool
	GRPCWrapper       grpc.Wrapper
}

type GenericMessage struct {
	Message string `json:"message"`
	Received time.Time `json:"received"`
	Sum [16]byte   `json:"sum"`
}
