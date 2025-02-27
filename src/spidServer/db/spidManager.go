package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
	"time"
)

func (m *Manager) GetSpidsFromFile() *entities.Spids {
	log.Print("Reading spids.")
	spidsFromFile, err := m.FileManager.ReadFile(DefaultSpidsLocation)
	if err != nil {
		log.Printf("Failed to read spids from file: %s", err)
		return nil
	}
	spids, err := entities.UnmarshalSpids(spidsFromFile)
	if err != nil {
		log.Printf("Failed to parse users from file: %s", err)
		return nil
	}
	if spids.Spids == nil {
		spids = entities.NewSpids()
	}
	log.Printf("Read spids: %s.", spids)
	return spids
}

func (m *Manager) WriteSpidsToFile() {
	marshaledSpids, err := entities.MarshalSpids(m.Spids)
	if err != nil {
		log.Printf("Failed to write spids to file: %s", err)
		return
	}

	log.Printf("Making backup of spids file...")
	src := m.FileManager.GetAbsolutePath(DefaultSpidsLocation)
	dst := src + ".bk"
	err = os.Rename(src, dst)
	eh.HandleFatal(err)

	log.Print("Writing spids.")
	err = m.FileManager.WriteToFile(DefaultSpidsLocation, marshaledSpids)
	if err != nil {
		log.Printf("Failed to write spids to file: %s", err)
		return
	}
	log.Print("Finished writing spids.")
}

func (m *Manager) QuerySpid(spidID uuid.UUID) (*entities.Spid, error) {
	log.Printf("Querying spid with ID %s.", spidID)
	s, ok := m.Spids.Spids[spidID]
	if !ok {
		err := fmt.Errorf("spid with ID %s not found", spidID)
		return nil, err
	}
	log.Printf("Spid found: %s", s)
	return s, nil
}

func (m *Manager) RegisterSpid(spid *entities.Spid) error {
	_, err := m.QuerySpid(spid.ID)
	if err == nil {
		return fmt.Errorf("spid with id %s already exists", spid.ID)
	}
	log.Printf("Registering spid: %s.", spid)
	m.Spids.Spids[spid.ID] = spid
	log.Print("Spid registered.")
	return m.logWriteAction(WriteAction{
		Location:   Local,
		EntityType: Spid,
		Type:       Register,
		SpidEntity:     spid,
	})
}

func (m *Manager) UpdateSpid(spid *entities.Spid) error {
	_, err := m.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating spid: %s.", spid)
	spid.LastUpdated = time.Now().Unix()
	m.Spids.Spids[spid.ID] = spid
	log.Print("Spid updated.")
	return m.logWriteAction(WriteAction{
		Location:   Local,
		EntityType: Spid,
		Type:       Update,
		SpidEntity:     spid,
	})
}

func (m *Manager) DeleteSpid(spidID uuid.UUID) error {
	spid, err := m.QuerySpid(spidID)
	if err != nil {
		return err
	}
	log.Printf("Deleting spid: %s.", spid)
	delete(m.Spids.Spids, spidID)
	log.Print("Spid deleted.")
	return m.logWriteAction(WriteAction{
		Location:   Local,
		EntityType: Spid,
		Type:       Delete,
		SpidEntity:     spid,
	})
}
