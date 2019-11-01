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
	log.Printf("Logging write action: %s", action)
	jsonWriteAction, err := json.Marshal(action)
	if err != nil {
		return err
	}
	log.Printf("Writing `%s`", jsonWriteAction)
	return m.DirtyLogger.Output(2, string(jsonWriteAction))
}

func (m *Manager) recoverWriteAction(data string) (writeAction WriteAction, err error) {
	log.Printf("Trying to recover write action from `%s`", data)
	err = json.Unmarshal([]byte(data), &writeAction)
	log.Printf("Recover write action: %s", writeAction)
	return writeAction, err
}

func (m *Manager) processWriteAction(action WriteAction) error {
	switch action.Location {
	default:
		return fmt.Errorf("invalid write action location `%s`", action.Location)
	case Local:
		return m.processLocalWriteAction(action)
	case Remote:
		return m.processRemoteWriteAction(action)
	}
}

func (m *Manager) processLocalWriteAction(action WriteAction) error {
	switch action.EntityType {
	default:
		return fmt.Errorf("invalid write action entity type `%s`", action.EntityType)
	case Spid:
		var handler func(*entities.Spid) error
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Register:
			handler = m.RegisterSpid
		case Update:
			handler = m.UpdateSpid
		case Delete:
			handler = m.DeleteSpid
		}
		spid := action.Entity.(entities.Spid)
		return handler(&spid)
	case User:
		var handler func(*entities.User) error
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Register:
			handler = m.RegisterUser
		case Update:
			handler = m.UpdateUser
		case Delete:
			handler = m.DeleteUser
		}
		user := action.Entity.(entities.User)
		return handler(&user)
	}
}

func (m *Manager) processRemoteWriteAction(action WriteAction) error {
	switch action.EntityType {
	default:
		return fmt.Errorf("invalid write action entity type `%s`", action.EntityType)
	case Spid:
		var handler func(*entities.Spid) error
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Add:
			handler = m.AddRemoteSpid
		case Update:
			handler = m.UpdateRemoteSpid
		case Remove:
			handler = m.RemoveRemoteSpid
		}
		spid := action.Entity.(entities.Spid)
		return handler(&spid)
	case User:
		var handler func(*entities.User) error
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Register:
			handler = m.AddRemoteUser
		case Update:
			handler = m.UpdateRemoteUser
		case Delete:
			handler = m.RemoveRemoteUser
		}
		user := action.Entity.(entities.User)
		return handler(&user)
	}
}
