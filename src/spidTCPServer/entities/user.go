package entities

import (
	"encoding/json"
	"github.com/google/uuid"
	"main/gps"
	"main/utils"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Name          string `json:"name"`
	Location      gps.GlobalPosition `json:"location"`
	CurrentSpidID string `json:"current_spid_id"`
}

type Users struct {
	Users []User `json:"users"`
}

func NewUser(name string) User {
	return User{uuid.New(), name, gps.NullPosition(), ""}
}

func (u User) Marshal() []byte {
	um, err := json.Marshal(u)
	utils.HandleFatal(err)
	return um
}

func Unmarshal(um []byte) User {
	var u User
	err := json.Unmarshal(um, &u)
	utils.HandleFatal(err)
	return u
}

func marshalUsers(users Users) []byte {
	m, err := json.Marshal(users)
	utils.HandleFatal(err)
	return m
}

// TODO: move to manager class
func (u User) Register(manager utils.DBManager) {
	var users Users
	err := json.Unmarshal(manager.GetUsersFromFile(), &users)
	utils.HandleFatal(err)
	users.Users = append(users.Users, u)

	manager.WriteUsersToFile(marshalUsers(users))
}
