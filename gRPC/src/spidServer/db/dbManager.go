package db

import (
	"github.com/google/uuid"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
	"spidServer/utils"
	"time"
)

const (
	Sep = string(os.PathSeparator)
	DefaultLogPath     = "requestHandling" + Sep + "request_logs.spd"
	DefaultDirtyRequestsPath = "requestHandling" + Sep + "dirty_requests.spd"
	DefaultMaxBufferedRequests = 100
	DefaultWriteToFilePeriod   = 5000*time.Millisecond

	DefaultUsersLocation = "db" + Sep + "users.spd"
	DefaultSpidsLocation = "db" + Sep + "spids.spd"
	DefaultServerIDLocation = "db" + Sep + "server_id.spd"
)

type Manager struct {
	FileManager     utils.FileManager
	Users           entities.Users
	Spids           entities.Spids
	LoggerDirty     *log.Logger
	WritingToMemory bool
	WritingToFile   bool
}

func NewManager(basePath string) Manager {
	pathDirty := basePath + Sep + DefaultDirtyRequestsPath
	dirtyLogFile, err := os.OpenFile(pathDirty, os.O_CREATE|os.O_RDWR, 0644)
	eh.HandleFatal(err)
	m := Manager{FileManager: utils.FileManager{BasePath: basePath}}
	m.LoggerDirty = log.New(dirtyLogFile, "", 0)
	m.LoadFromFile()
	go m.WriteToFilePeriodically(DefaultWriteToFilePeriod)
	return m
}

func (m *Manager) LoadFromFile() {
	m.Users = m.GetUsersFromFile()
	m.Spids = m.GetSpidsFromFile()
}

func (m *Manager) GetServerID() uuid.UUID {
	serverIDPath := DefaultServerIDLocation
	id, err := m.FileManager.ReadFile(serverIDPath)
	if err != nil {
		return uuid.Nil
	}
	uid, err := uuid.Parse(string(id))
	if err != nil {
		return uuid.Nil
	}
	return uid
}

func (m *Manager) WriteServerID(uid uuid.UUID) error {
	return m.FileManager.WriteToFile(DefaultServerIDLocation, []byte(uid.String()))
}

func (m *Manager) WriteToFilePeriodically(period time.Duration) {
	for {
		time.Sleep(period)
		for ; m.WritingToMemory; {}
		m.WritingToFile = true
		log.Print("Writing users and spids to file.")
		m.WriteUsersToFile()
		m.WriteSpidsToFile()

		log.Print("Truncating dirty log file...")
		pathDirty := m.FileManager.BasePath + Sep + DefaultDirtyRequestsPath
		dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		// Old *File for LoggerDirty will be garbage collected
		m.LoggerDirty = log.New(dirtyLogFile, "", 0)

		eh.HandleFatal(err)
		m.WritingToFile = false
	}
}
