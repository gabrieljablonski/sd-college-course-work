package entities

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"spidServer/gps"
	"spidServer/requestHandling/protoBuffers"
	"time"
)

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

type User struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	Position      gps.GlobalPosition `json:"position"`
	LastUpdated   int64              `json:"last_updated"`
	CurrentSpidID uuid.UUID          `json:"current_spid_id"`
}

func UserFromProtoBufferEntity(pbUser *protoBuffers.User) (user User, err error) {
	idU, err := uuid.Parse(pbUser.Id)
	if err != nil {
		return user, err
	}
	idS, err := uuid.Parse(pbUser.CurrentSpidID)
	if err != nil {
		return user, err
	}
	return User{
		ID:            idU,
		Name:          pbUser.Name,
		Position:      gps.FromProtoBufferEntity(pbUser.Position),
		LastUpdated:   pbUser.LastUpdated,
		CurrentSpidID: idS,
	}, nil
}

func (u User) ToString() string {
	s, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to convert user to string: %s", err)
		return ""
	}
	return string(s)
}

func (u User) ToProtoBufferEntity() *protoBuffers.User {
	return &protoBuffers.User{
		Id:            u.ID.String(),
		Name:          u.Name,
		Position:      u.Position.ToProtoBufferEntity(),
		LastUpdated:   u.LastUpdated,
		CurrentSpidID: u.CurrentSpidID.String(),
	}
}

func NewUser(name string, position gps.GlobalPosition) User {
	return User{
		ID:            uuid.New(),
		Name:          name,
		Position:      position,
		LastUpdated:   time.Now().Unix(),
		CurrentSpidID: uuid.Nil,
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

func (u *User) UpdatePosition(position gps.GlobalPosition) {
	u.Position = position
	u.LastUpdated = time.Now().Unix()
}
