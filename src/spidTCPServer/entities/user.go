package entities

import (
	"encoding/json"
	"github.com/google/uuid"
	eh "main/errorHandling"
	"main/gps"
	"time"
)

type User struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	OnRide        bool               `json:"on_ride"`
	Location      gps.GlobalPosition `json:"location"`
	LastUpdated   time.Time          `json:"last_updated"`
	CurrentSpidID string             `json:"current_spid_id"`
}

type Users struct {
	Users map[uuid.UUID]User `json:"users"`
}

func NewUser(name string) User {
	return User{
		uuid.New(),
		name,
		false,
		gps.NullPosition(),
		time.Unix(0,0),
		"",
	}
}

func (u User) Marshal() []byte {
	um, err := json.Marshal(u)
	eh.HandleFatal(err)
	return um
}

func UnmarshalUsers(marshaledUsers []byte) Users {
	var users Users
	err := json.Unmarshal(marshaledUsers, &users)
	eh.HandleFatal(err)
	return users
}

func MarshalUsers(users Users) []byte {
	m, err := json.MarshalIndent(users, "", "    ")
	eh.HandleFatal(err)
	return m
}

func (u User) UpdateLocation(position gps.GlobalPosition) {
	u.Location = position
	u.LastUpdated = time.Now()
}
