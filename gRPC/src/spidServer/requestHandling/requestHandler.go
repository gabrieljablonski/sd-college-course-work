package requestHandling

import (
	"log"
	"os"
	"spidServer/db"
	eh "spidServer/errorHandling"
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
	Manager           db.Manager
	LoggerDirty       *log.Logger
	WritingToFile     bool
	WritingToMemory   bool
}

func NewHandler(basePath string) Handler {
	pathDirty := basePath + string(os.PathSeparator) + DefaultDirtyRequestsPath

	dirtyLogFile, err := os.OpenFile(pathDirty, os.O_CREATE|os.O_RDWR, 0644)
	eh.HandleFatal(err)
	h := Handler{
		Manager: db.NewManager(basePath),
	}
	h.Manager.LoadFromFile()
	h.processDirtyRequests(dirtyLogFile)
	h.LoggerDirty = log.New(dirtyLogFile, "", 0)
	go h.WriteToFilePeriodically(DefaultWriteToFilePeriod)
	return h
}

func (h *Handler) processDirtyRequests(dirtyLogFile *os.File) {
	//scanner := bufio.NewScanner(dirtyLogFile)
	//log.Print("Processing requests from last session...")
	//var requestMessage GenericMessage
	//for scanner.Scan() {
	//	err := json.Unmarshal([]byte(scanner.Text()), &requestMessage)
	//	eh.HandleFatal(err)
	//	log.Printf("Processing request: `%s`", requestMessage.Message)
	//	h.processRequest(requestMessage.Message)
	//}
	//err := scanner.Err()
	//eh.HandleFatal(err)
	//
	//h.Manager.WriteUsersToFile()
	//h.Manager.WriteSpidsToFile()
	//err = dirtyLogFile.Truncate(0)
	//eh.HandleFatal(err)
	//_, err = dirtyLogFile.Seek(0, 0)
	//eh.HandleFatal(err)
	//log.Printf("Memory and file should be up to date to last session.")
}

func (h *Handler) WriteToFilePeriodically(period time.Duration) {
	for {
		time.Sleep(period)
		for ; h.WritingToMemory; {}
		h.WritingToFile = true
		log.Print("Writing users and spids to file.")
		h.Manager.WriteUsersToFile()
		h.Manager.WriteSpidsToFile()

		log.Print("Truncating dirty log file...")
		pathDirty := h.Manager.FileManager.BasePath + string(os.PathSeparator) + DefaultDirtyRequestsPath
		dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		// Old *File for LoggerDirty will be garbage collected
		h.LoggerDirty = log.New(dirtyLogFile, "", 0)

		eh.HandleFatal(err)
		h.WritingToFile = false
	}
}
