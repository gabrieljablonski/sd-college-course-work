package db

import (
	"bufio"
	"github.com/google/uuid"
	"io"
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
	DirtyLogger     *log.Logger
	WritingToMemory bool
	WritingToFile   bool
}

func NewManager(basePath string) Manager {
	m := Manager{FileManager: utils.FileManager{BasePath: basePath}}
	m.loadFromFile()
	pathDirty := m.FileManager.GetAbsolutePath(DefaultDirtyRequestsPath)
	log.Print("Creating new log file...")
	dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	eh.HandleFatal(err)
	m.DirtyLogger = log.New(dirtyLogFile, "", 0)
	go m.WriteToFilePeriodically(DefaultWriteToFilePeriod)
	return m
}

func (m *Manager) recoverFromSavedLogs() {
	log.Print("Opening log file to recover data...")
	pathDirty := m.FileManager.GetAbsolutePath(DefaultDirtyRequestsPath)
	dirtyLogFile, err := os.Open(pathDirty)
	if os.IsNotExist(err) {
		log.Print("Log file does not exist.")
		return
	}
	defer func() {
		_ = dirtyLogFile.Close()
	}()
	reader := bufio.NewReader(dirtyLogFile)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Print("Finished recovering from saved logs.")
			break
		} else if err != nil {
			log.Fatalf("Failed to recover from saved logs: %s", err)
		}
		writeAction, err := m.recoverWriteAction(line)
		err = m.processWriteAction(writeAction)
		if err != nil {
			log.Fatalf("Failed to process write action: %s", err)
		}
	}
}

func (m *Manager) loadFromFile() {
	log.Print("Recovering previous server state from files...")
	m.Users = m.GetUsersFromFile()
	m.Spids = m.GetSpidsFromFile()
	m.RemoteUsers = m.GetRemoteUsersFromFile()
	m.RemoteSpids = m.GetRemoteSpidsFromFile()
	m.recoverFromSavedLogs()
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
		pathDirty := m.FileManager.GetAbsolutePath(DefaultDirtyRequestsPath)
		dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		// Old *File for DirtyLogger will be garbage collected
		m.DirtyLogger = log.New(dirtyLogFile, "", 0)

		eh.HandleFatal(err)
		m.WritingToFile = false
	}
}

func (m *Manager) GetServerID() uuid.UUID {
	log.Print("Recovering server ID from file...")
	serverIDPath := DefaultServerIDLocation
	id, err := m.FileManager.ReadFile(serverIDPath)
	if err != nil {
		log.Printf("Unable to read file: %s", err)
		return uuid.Nil
	}
	uid, err := uuid.Parse(string(id))
	if err != nil {
		log.Printf("Failed to parse id `%s`: %s", id, err)
		return uuid.Nil
	}
	log.Printf("Recovered %s", uid.String())
	return uid
}

func (m *Manager) WriteServerID(uid uuid.UUID) error {
	return m.FileManager.WriteToFile(DefaultServerIDLocation, []byte(uid.String()))
}
