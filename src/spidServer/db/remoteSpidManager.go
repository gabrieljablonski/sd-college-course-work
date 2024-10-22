package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
)

func (m *Manager) GetRemoteSpidsFromFile() *entities.Spids {
	log.Print("Reading remote spids.")
	spidsFromFile, err := m.FileManager.ReadFile(DefaultRemoteSpidsLocation)
	if err != nil {
		log.Printf("Failed to read remote spids from file: %s", err)
		return nil
	}
	spids, err := entities.UnmarshalSpids(spidsFromFile)
	if err != nil {
		log.Printf("Failed to parse remote spids from file: %s", err)
		return nil
	}
	if spids.Spids == nil {
		spids = entities.NewSpids()
	}
	log.Printf("Read remote spids: %s.", spids)
	return spids
}

func (m *Manager) WriteRemoteSpidsToFile() {
	marshaledSpids, err := entities.MarshalSpids(m.RemoteSpids)
	if err != nil {
		log.Printf("Failed to write remote spids to file: %s", err)
		return
	}

	log.Printf("Making backup of remote spids file...")
	src := m.FileManager.GetAbsolutePath(DefaultRemoteSpidsLocation)
	dst := src + ".bk"
	err = os.Rename(src, dst)
	eh.HandleFatal(err)

	log.Print("Writing remote spids.")
	err = m.FileManager.WriteToFile(DefaultRemoteSpidsLocation, marshaledSpids)
	if err != nil {
		log.Printf("Failed to write remote spids to file: %s", err)
		return
	}
	log.Print("Finished writing remote spids.")
}

func (m *Manager) GetRemoteSpids() *entities.Spids {
	log.Print("Querying remote spids.")
	return m.RemoteSpids
}

func (m *Manager) QueryRemoteSpid(spidID uuid.UUID) (*entities.Spid, error) {
	log.Printf("Querying remote spid with ID %s.", spidID)
	s, ok := m.RemoteSpids.Spids[spidID]
	if !ok {
		err := fmt.Errorf("remote spid with ID %s not found", spidID)
		return nil, err
	}
	log.Printf("Remote spid found: %s", s)
	return s, nil
}

func (m *Manager) AddRemoteSpid(spid *entities.Spid) error {
	log.Printf("Adding remote spid: %s.", spid)
	_, err := m.QueryRemoteSpid(spid.ID)
	if err == nil {
		return fmt.Errorf("spid with id %s already exists", spid.ID)
	}
	m.RemoteSpids.Spids[spid.ID] = spid
	log.Print("Remote spid added.")
	return m.logWriteAction(WriteAction{
		Location:   Remote,
		EntityType: Spid,
		Type:       Add,
		SpidEntity: spid,
	})
}

func (m *Manager) UpdateRemoteSpid(spid *entities.Spid) error {
	log.Printf("Updating remote spid: %s.", spid)
	_, err := m.QueryRemoteSpid(spid.ID)
	if err != nil {
		return err
	}
	m.RemoteSpids.Spids[spid.ID] = spid
	log.Print("Remote spid updated.")
	return m.logWriteAction(WriteAction{
		Location:   Remote,
		EntityType: Spid,
		Type:       Update,
		SpidEntity: spid,
	})
}

func (m *Manager) RemoveRemoteSpid(spidID uuid.UUID) error {
	spid, err := m.QueryRemoteSpid(spidID)
	if err != nil {
		return err
	}
	log.Printf("Removing remote spid: %s.", spid)
	delete(m.RemoteSpids.Spids, spid.ID)
	log.Print("Remote spid removed.")
	return m.logWriteAction(WriteAction{
		Location:   Remote,
		EntityType: Spid,
		Type:       Remove,
		SpidEntity: spid,
	})
}
