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
	SpidEntity *entities.Spid        `json:"spid"`
	UserEntity *entities.User		 `json:"user"`
}

//func (w WriteAction) String() string {
//	return fmt.Sprintf("%#v", w)
//}

func (w WriteAction) Json() (string, error) {
	s, err := json.Marshal(w)
	return string(s), err
}

func (m *Manager) logWriteAction(action WriteAction) error {
	if m.DirtyLogger == nil {
		return nil
	}
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
	log.Printf("Recovered write action: %s", writeAction)
	return writeAction, err
}

func (m *Manager) ProcessWriteAction(action WriteAction) error {
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
		spid := action.SpidEntity
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Register:
			return m.RegisterSpid(spid)
		case Update:
			return m.UpdateSpid(spid)
		case Delete:
			return m.DeleteSpid(spid.ID)
		}
	case User:
		user := action.UserEntity
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Register:
			return m.RegisterUser(user)
		case Update:
			return m.UpdateUser(user)
		case Delete:
			return m.DeleteUser(user.ID)
		}
	}
}

func (m *Manager) processRemoteWriteAction(action WriteAction) error {
	switch action.EntityType {
	default:
		return fmt.Errorf("invalid write action entity type `%s`", action.EntityType)
	case Spid:
		spid := action.SpidEntity
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Add:
			return m.AddRemoteSpid(spid)
		case Update:
			return m.UpdateRemoteSpid(spid)
		case Remove:
			return m.RemoveRemoteSpid(spid.ID)
		}
	case User:
		user := action.UserEntity
		switch action.Type {
		default:
			return fmt.Errorf("invalid write action type `%s`", action.Type)
		case Add:
			return m.AddRemoteUser(user)
		case Update:
			return m.UpdateRemoteUser(user)
		case Delete:
			return m.RemoveRemoteUser(user.ID)
		}
	}
}
