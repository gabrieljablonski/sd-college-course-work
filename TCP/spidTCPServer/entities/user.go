package entities

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"main/gps"
	"time"
)

type User struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	Location      gps.GlobalPosition `json:"location"`
	LastUpdated   time.Time          `json:"last_updated"`
	CurrentSpidID uuid.UUID          `json:"current_spid_id"`
}

func (u User) ToString() string {
	s, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to convert user to string: %s", err)
		return ""
	}
	return string(s)
}

type Users struct {
	Users map[uuid.UUID]User `json:"users"`
}

func (u Users) ToString() string {
	s, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to convert users to string: %s", err)
		return ""
	}
	return string(s)
}

func NewUser(name string) User {
	return User{
		uuid.New(),
		name,
		gps.NullPosition(),
		time.Unix(0,0),
		uuid.Nil,
	}
}

func (u User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func MarshalUsers(users Users) ([]byte, error) {
	return json.MarshalIndent(users, "", "    ")
}

func UnmarshalUsers(marshaledUsers []byte) (Users, error) {
	var users Users
	err := json.Unmarshal(marshaledUsers, &users)
	return users, err
}

func (u *User) UpdateLocation(position gps.GlobalPosition) {
	u.Location = position
	u.LastUpdated = time.Now()
}