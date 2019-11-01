package db

import (
	"encoding/json"
)

const (
	Spid = iota
	User

	RegisterLocal
	UpdateLocal
	DeleteLocal

	AddRemote
	UpdateRemote
	RemoveRemote
)

type WriteActionEntityType int

func (e WriteActionEntityType) String() string {
	entities := map[WriteActionEntityType]string{
		User: "USER",
		Spid: "SPID",
	}
	return entities[e]
}

type WriteActionType int

func (a WriteActionType) String() string {
	types := map[WriteActionType]string{
		RegisterLocal: "REGISTER",
		UpdateLocal:   "UPDATE",
		DeleteLocal:   "DELETE",
		AddRemote:     "ADD REMOTE",
		UpdateRemote:  "UPDATE REMOTE",
		RemoveRemote:  "REMOVE REMOTE",
	}
	return types[a]
}

type WriteAction struct {
	Type   WriteActionType           `json:"type"`
	EntityType WriteActionEntityType `json:"entity_type"`
	Entity interface{}               `json:"entity"`
}

func (m *Manager) logWriteAction(action WriteAction) error {
	jsonWriteAction, err := json.Marshal(action)
	if err != nil {
		return err
	}
	return m.DirtyLogger.Output(2, string(jsonWriteAction))
}

func (m *Manager) recoverWriteAction(log []byte) (writeAction WriteAction, err error) {
	err = json.Unmarshal(log, &writeAction)
	return writeAction, err
}
