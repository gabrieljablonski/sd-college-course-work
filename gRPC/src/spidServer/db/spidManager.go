package db

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"spidServer/entities"
	eh "spidServer/errorHandling"
)

func (m *Manager) GetSpidsFromFile() entities.Spids {
	log.Print("Reading spids.")
	spidsFromFile, err := m.FileManager.ReadFile(DefaultSpidsLocation)
	if err != nil {
		log.Printf("Failed to read spids from file: %s", err)
		return entities.Spids{}
	}
	spids, err := entities.UnmarshalSpids(spidsFromFile)
	if err != nil {
		log.Printf("Failed to parse users from file: %s", err)
		return entities.Spids{}
	}
	if spids.Spids == nil {
		spids.Spids = map[uuid.UUID]entities.Spid{}
	}
	log.Printf("Read spids: %s.", spids.ToString())
	return spids
}

func (m *Manager) WriteSpidsToFile() {
	marshaledSpids, err := entities.MarshalSpids(m.Spids)
	if err != nil {
		log.Printf("Failed to write spids to file: %s", err)
		return
	}

	log.Printf("Making copy of spids file...")
	src := m.FileManager.BasePath + string(os.PathSeparator) + DefaultSpidsLocation
	dst := src + ".bk"
	source, err := os.Open(src)
	if err != nil {
		eh.HandleFatal(err)
	}
	defer source.Close()
	_ = os.Remove(dst)
	destination, err := os.Create(dst)
	if err != nil {
		eh.HandleFatal(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	eh.HandleFatal(err)
	
	log.Printf("Writing spids: %s", m.Spids.ToString())
	err = m.FileManager.WriteToFile(DefaultSpidsLocation, marshaledSpids)
	if err != nil {
		log.Printf("Failed to write spids to file: %s", err)
		return
	}
	log.Print("Finished writing spids.")
}

func (m *Manager) QuerySpid(spidID uuid.UUID) (entities.Spid, error) {
	log.Printf("Querying spid with ID %s.", spidID)
	_, ok := m.Spids.Spids[spidID]
	if !ok {
		err := fmt.Errorf("spid with ID %s not found", spidID)
		return entities.Spid{}, err
	}
	log.Printf("Spid found: %s", m.Spids.Spids[spidID])
	return m.Spids.Spids[spidID], nil
}

func (m *Manager) RegisterSpid(spid entities.Spid) error {
	log.Printf("Registering spid: %s.", spid.ToString())
	m.Spids.Spids[spid.ID] = spid
	log.Print("Spid registered.")
	return nil
}

func (m *Manager) UpdateSpid(spid entities.Spid) error {
	_, err := m.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating spid: %s.", spid.ToString())
	m.Spids.Spids[spid.ID] = spid
	log.Print("Spid updated.")
	return nil
}

func (m *Manager) DeleteSpid(spid entities.Spid) error {
	_, err := m.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleting spid: %s.", spid.ToString())
	delete(m.Spids.Spids, spid.ID)
	log.Print("Spid updated.")
	return nil
}
