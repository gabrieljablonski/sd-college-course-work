package db

import (
	"encoding/json"
	"fmt"
	"log"
	"spidServer/entities"
)

const (
	Spid = iota
	User

	Local
	Remote

	Register
	Update
	Delete
	Add
	Remove
)

type WriteActionType int

func (a WriteActionType) String() string {
	types := map[WriteActionType]string{
		Register: "REGISTER",
		Update:   "UPDATE",
		Delete:   "DELETE",
		Add:      "ADD",
		Remove:   "REMOVE",
	}
	return types[a]
}

type WriteActionLocation int

func (e WriteActionLocation) String() string {
	locations := map[WriteActionLocation]string{
		Local: "LOCAL",
		Remote: "REMOTE",
	}
	return locations[e]
}

type WriteActionEntityType int

func (e WriteActionEntityType) String() string {
	entityTypes := map[WriteActionEntityType]string{
		User: "USER",
		Spid: "SPID",
	}
	return entityTypes[e]
}

type WriteAction struct {
	Location   WriteActionLocation   `json:"location"`
	EntityType WriteActionEntityType `json:"entity_type"`
	Type       WriteActionType       `json:"type"`
	Entity     interface{}           `json:"entity"`
}

func (w WriteAction) String() string {
	return fmt.Sprintf("WriteAction{Location: %s, EntityType: %s, Type: %s, Entity: %s}",
		               w.Location, w.EntityType, w.Type, w.Entity)
}

func (m *Manager) logWriteAction(action WriteAction) error {
	jsonWriteAction, err := json.Marshal(action)
	if err != nil {
		return err
	}
	return m.DirtyLogger.Output(2, string(jsonWriteAction))
}

func (m *Manager) recoverWriteAction(data string) (writeAction WriteAction, err error) {
	log.Printf("Trying to recover write action from `%s`", data)
	err = json.Unmarshal([]byte(data), &writeAction)
	log.Printf("Recover write action: %s", writeAction)
	return writeAction, err
}
