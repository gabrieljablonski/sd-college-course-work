package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
)

func (m *Manager) GetRemoteSpidsFromFile() entities.Spids {
	log.Print("Reading remote spids.")
	spidsFromFile, err := m.FileManager.ReadFile(DefaultRemoteSpidsLocation)
	if err != nil {
		log.Printf("Failed to read remote spids from file: %s", err)
		return entities.Spids{}
	}
	spids, err := entities.UnmarshalSpids(spidsFromFile)
	if err != nil {
		log.Printf("Failed to parse remote spids from file: %s", err)
		return entities.Spids{}
	}
	if spids.Spids == nil {
		spids.Spids = map[uuid.UUID]entities.Spid{}
	}
	log.Printf("Read remote spids: %s.", spids.ToString())
	return spids
}

func (m *Manager) WriteRemoteSpidsToFile() {
	marshaledSpids, err := entities.MarshalSpids(m.RemoteSpids)
	if err != nil {
		log.Printf("Failed to write remote spids to file: %s", err)
		return
	}

	log.Printf("Making backup of remote spids file...")
	src := m.FileManager.BasePath + string(os.PathSeparator) + DefaultRemoteSpidsLocation
	dst := src + ".bk"
	err = os.Rename(src, dst)
	eh.HandleFatal(err)

	log.Printf("Writing remote spids: %s", m.RemoteSpids.ToString())
	err = m.FileManager.WriteToFile(DefaultRemoteSpidsLocation, marshaledSpids)
	if err != nil {
		log.Printf("Failed to write remote spids to file: %s", err)
		return
	}
	log.Print("Finished writing remote spids.")
}

func (m *Manager) GetRemoteSpids() entities.Spids {
	log.Print("Querying remote spids.")
	return m.RemoteSpids
}

func (m *Manager) QueryRemoteSpid(spidID uuid.UUID) (entities.Spid, error) {
	log.Printf("Querying remote spid with ID %s.", spidID)
	s, ok := m.RemoteSpids.Spids[spidID]
	if !ok {
		err := fmt.Errorf("remote spid with ID %s not found", spidID)
		return entities.Spid{}, err
	}
	log.Printf("Remote spid found: %s", s.ToString())
	return s, nil
}

func (m *Manager) AddRemoteSpid(spid entities.Spid) error {
	log.Printf("Adding remote spid: %s.", spid.ToString())
	m.RemoteSpids.Spids[spid.ID] = spid
	log.Print("Remote spid added.")
	return nil
}

func (m *Manager) UpdateRemoteSpid(spid entities.Spid) error {
	_, err := m.QueryRemoteSpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating remote spid: %s.", spid.ToString())
	m.RemoteSpids.Spids[spid.ID] = spid
	log.Print("Remote spid updated.")
	return nil
}

func (m *Manager) RemoveRemoteSpid(spid entities.Spid) error {
	_, err := m.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Removing remote spid: %s.", spid.ToString())
	delete(m.Spids.Spids, spid.ID)
	log.Print("Remote spid removed.")
	return nil
}
