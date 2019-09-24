package entities

import (
	"encoding/json"
	"spidTCPServer/external/uuid"
	"spidTCPServer/gps"
	"spidTCPServer/utils"
)

type User struct {
	ID            uuid.UUID
	Name          string
	Location      gps.GlobalPosition
	CurrentSpidID string
}

func New(name string) User {
	return User{uuid.New(), name, gps.NullPosition(), nil}
}

func (u User) marshal() []byte {
	um, err := json.Marshal(u)
	utils.HandleFatal(err)
	return um
}

func unmarshal(um []byte) User {
	var u User
	err := json.Unmarshal(um, &u)
	utils.HandleFatal(err)
	return u
}
