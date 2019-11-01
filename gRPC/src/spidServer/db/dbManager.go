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
	DefaultMaxBufferedRequests = 100
	DefaultWriteToFilePeriod   = 5000000*time.Millisecond

	Sep                        = string(os.PathSeparator)
	BaseDataPath               = "data" + Sep
	BaseLogsPath               = BaseDataPath + "logs" + Sep
	BaseStatePath              = BaseDataPath + "state" + Sep
	DefaultDirtyRequestsPath   = BaseLogsPath + "dirty_requests.spd"
	DefaultUsersLocation 	   = BaseStatePath + "users.spd"
	DefaultRemoteUsersLocation = BaseStatePath + "users_remote.spd"
	DefaultSpidsLocation       = BaseStatePath + "spids.spd"
	DefaultRemoteSpidsLocation = BaseStatePath + "spids_remote.spd"
	DefaultServerIDLocation    = BaseDataPath + "server_id.spd"
)

type Manager struct {
	FileManager     utils.FileManager
	Users           *entities.Users
	Spids           *entities.Spids
	RemoteUsers     *entities.Users
	RemoteSpids     *entities.Spids
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
	m.loadFromFile()
	go m.WriteToFilePeriodically(DefaultWriteToFilePeriod)
	return m
}

func (m *Manager) loadFromFile() {
	m.Users = m.GetUsersFromFile()
	m.Spids = m.GetSpidsFromFile()
	m.RemoteUsers = m.GetRemoteUsersFromFile()
	m.RemoteSpids = m.GetRemoteSpidsFromFile()
}

func (m *Manager) WriteToFilePeriodically(period time.Duration) {
	for {
		time.Sleep(period)
		for ; m.WritingToMemory; {}
		m.WritingToFile = true
		log.Print("Writing users and spids to file.")
		m.WriteUsersToFile()
		m.WriteSpidsToFile()
		m.WriteRemoteUsersToFile()
		m.WriteRemoteSpidsToFile()

		log.Print("Truncating dirty log file...")
		pathDirty := m.FileManager.BasePath + Sep + DefaultDirtyRequestsPath
		dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		// Old *File for LoggerDirty will be garbage collected
		m.LoggerDirty = log.New(dirtyLogFile, "", 0)

		eh.HandleFatal(err)
		m.WritingToFile = false
	}
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
