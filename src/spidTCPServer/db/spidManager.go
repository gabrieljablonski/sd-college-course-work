package db

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"main/entities"
)

// TODO: error handling for file operations

func (d Manager) GetSpidsFromFile() entities.Spids {
	log.Print("Reading spids.")
	spids := entities.UnmarshalSpids(d.FM.ReadFile(SpidsDefaultLocation))
	log.Printf("Read spids: %s.", spids)
	return spids
}

func (d Manager) WriteSpidsToFile(spids []byte) {
	log.Printf("Writing spids: %s", string(spids))
	d.FM.WriteToFile(SpidsDefaultLocation, spids)
	log.Print("Finished writing spids.")
}

func (d Manager) writeSpid(spid entities.Spid) {
	spids := d.GetSpidsFromFile()
	spids.Spids[spid.ID] = spid
	log.Printf("Writing spid: %s.", spid)
	d.WriteSpidsToFile(entities.MarshalSpids(spids))
	log.Print("Spid written.")
}

func (d Manager) QuerySpid(spidID uuid.UUID) (entities.Spid, error) {
	spids := d.GetSpidsFromFile()
	log.Printf("Querying spid with ID %s.", spidID)
	_, ok := spids.Spids[spidID]
	if !ok {
		err := fmt.Errorf("spid with ID %s not found", spidID)
		return entities.Spid{}, err
	}
	log.Printf("Spid found: %s", spids.Spids[spidID])
	return spids.Spids[spidID], nil
}

func (d Manager) RegisterSpid(spid entities.Spid) error {
	_, err := d.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Registering spid: %s.", spid)
	d.writeSpid(spid)
	log.Print("Spid registered.")
	return nil
}

func (d Manager) UpdateSpid(spid entities.Spid) error {
	_, err := d.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Updating spid: %s.", spid)
	d.writeSpid(spid)
	log.Print("Spid updated.")
	return nil
}

func (d Manager) DeleteSpid(spid entities.Spid) error {
	_, err := d.QuerySpid(spid.ID)
	if err != nil {
		return err
	}
	log.Printf("Deleting spid: %s.", spid)
	spids := d.GetSpidsFromFile()
	delete(spids.Spids, spid.ID)
	d.WriteSpidsToFile(entities.MarshalSpids(spids))
	log.Print("Spid updated.")
	return nil
}
