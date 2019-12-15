package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
	"spidServer/utils"
	"time"
)

const (
	DefaultMaxBufferedRequests = 100
	DefaultWriteToFilePeriod   = 5000*time.Hour

	Sep                        = string(os.PathSeparator)
	BaseDataPath               = "data" + Sep
	DefaultServerIDLocation    = BaseDataPath + "server_id.spd"
	DefaultIPMapLocation       = BaseDataPath + "ip_map.spd"
	BaseStatePath              = BaseDataPath + "state" + Sep
	DefaultDirtyRequestsPath   = BaseStatePath + "dirty_requests.spd"
	DefaultUsersLocation 	   = BaseStatePath + "users.spd"
	DefaultRemoteUsersLocation = BaseStatePath + "users_remote.spd"
	DefaultSpidsLocation       = BaseStatePath + "spids.spd"
	DefaultRemoteSpidsLocation = BaseStatePath + "spids_remote.spd"
)

type Manager struct {
	FileManager      utils.FileManager
	Users            *entities.Users
	Spids            *entities.Spids
	RemoteUsers      *entities.Users
	RemoteSpids      *entities.Spids
	DirtyLogger      *log.Logger
}

func NewManager(basePath string) Manager {
	m := Manager{FileManager: utils.FileManager{BasePath: basePath}}
	m.loadFromFile()
	pathDirty := m.FileManager.GetAbsolutePath(DefaultDirtyRequestsPath)
	log.Print("Creating new log file...")
	dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	eh.HandleFatal(err)
	m.DirtyLogger = log.New(dirtyLogFile, "", 0)
	go m.writeToFilePeriodically(DefaultWriteToFilePeriod)
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
		err = m.ProcessWriteAction(writeAction)
		if err != nil {
			log.Fatalf("Failed to process write action: %s", err)
		}
	}
	m.writeEntitiesToFile()
}

func (m *Manager) createStateFiles() {
	statePath := m.FileManager.GetAbsolutePath(BaseStatePath)
	if err := os.MkdirAll(statePath, 0777); err != nil {
		if !os.IsExist(err){
			log.Fatalf("Failed to create state directory: %v", err)
		}
	}
	fileContentMap := map[string][]byte{
		DefaultUsersLocation: []byte("{\"users\":{}}"),
		DefaultRemoteUsersLocation: []byte("{\"users\":{}}"),
		DefaultSpidsLocation: []byte("{\"spids\":{}}"),
		DefaultRemoteSpidsLocation: []byte("{\"spids\":{}}"),
	}

	for f, c := range fileContentMap {
		if err := ioutil.WriteFile(m.FileManager.GetAbsolutePath(f), c, 0777); err != nil {
			log.Fatalf("Failed to create state directory: %v", err)
		}
	}
}

func (m *Manager) loadFromFile() {
	log.Print("Recovering previous server state from files...")
	m.createStateFiles()
	m.Users = m.GetUsersFromFile()
	m.Spids = m.GetSpidsFromFile()
	m.RemoteUsers = m.GetRemoteUsersFromFile()
	m.RemoteSpids = m.GetRemoteSpidsFromFile()
	//m.recoverFromSavedLogs()
}

func (m *Manager) writeEntitiesToFile() {
	log.Print("Writing users and spids to file.")
	m.WriteUsersToFile()
	m.WriteSpidsToFile()
	m.WriteRemoteUsersToFile()
	m.WriteRemoteSpidsToFile()
	log.Print("Truncating dirty log file...")
	m.truncateLogFile()
}

func (m *Manager) truncateLogFile() {
	pathDirty := m.FileManager.GetAbsolutePath(DefaultDirtyRequestsPath)
	dirtyLogFile, err := os.OpenFile(pathDirty, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	// Old *File for DirtyLogger will be garbage collected
	m.DirtyLogger = log.New(dirtyLogFile, "", 0)
	eh.HandleFatal(err)
}

func (m *Manager) writeToFilePeriodically(period time.Duration) {
	for {
		time.Sleep(period)
		m.writeEntitiesToFile()
	}
}

func (m *Manager) GetServerIDFromFile() uuid.UUID {
	log.Print("Recovering server ID from file...")
	serverIDPath := DefaultServerIDLocation
	id, err := m.FileManager.ReadFile(serverIDPath)
	if err != nil {
		log.Printf("%s", err)
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

func (m *Manager) WriteServerIDToFile(uid uuid.UUID) error {
	log.Printf("Saving server ID `%s` to file", uid.String())
	return m.FileManager.WriteToFile(DefaultServerIDLocation, []byte(uid.String()))
}

func (m *Manager) GetIPMapFromFile() (map[int][]utils.IP, error) {
	log.Print("Recovering IP map from file...")
	ipMapPath := DefaultIPMapLocation
	content, err := m.FileManager.ReadFile(ipMapPath)
	if err != nil {
		return nil, err
	}
	var ipMap map[int][]utils.IP
	err = json.Unmarshal(content, &ipMap)
	if err != nil {
		return nil, err
	}
	if len(ipMap) == 0 {
		return nil, fmt.Errorf("ip map is empty")
	}
	log.Printf("Recovered %v", ipMap)
	return ipMap, nil
}

func (m *Manager) WriteIPMapToFile(ipMap map[int][]utils.IP) error {
	log.Print("Saving IP map to file...")
	ipMapString, err := json.Marshal(ipMap)
	if err != nil {
		return err
	}
	return m.FileManager.WriteToFile(DefaultIPMapLocation, ipMapString)
}
