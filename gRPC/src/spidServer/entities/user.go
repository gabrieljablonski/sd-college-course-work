package entities

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"spidServer/gps"
	"spidServer/requestHandling/protoBuffers"
	"time"
)

type Users struct {
	Users map[uuid.UUID]*User `json:"users"`
}

func NewUsers() *Users {
	return &Users{
		Users: map[uuid.UUID]*User{},
	}
}

func (u Users) String() string {
	s, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to convert users to string: %s", err)
		return ""
	}
	return string(s)
}

type User struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"name"`
	Position      gps.GlobalPosition  `json:"position"`
	LastUpdated   int64               `json:"last_updated"`
	CurrentSpidID uuid.UUID           `json:"current_spid_id"`
}

func UserFromProtoBufferEntity(pbUser *protoBuffers.User) (user *User, err error) {
	idU, err := uuid.Parse(pbUser.Id)
	if err != nil {
		return user, err
	}
	idS, err := uuid.Parse(pbUser.CurrentSpidID)
	if err != nil {
		return user, err
	}
	position, err := gps.FromProtoBufferEntity(pbUser.Position)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:            idU,
		Name:          pbUser.Name,
		Position:      position,
		LastUpdated:   pbUser.LastUpdated,
		CurrentSpidID: idS,
	}, nil
}

func (u User) String() string {
	s, err := json.Marshal(u)
	if err != nil {
		log.Printf("Failed to convert user to string: %s", err)
		return ""
	}
	return string(s)
}

func (u User) ToProtoBufferEntity() (*protoBuffers.User, error) {
	position, err := u.Position.ToProtoBufferEntity()
	if err != nil {
		return nil, err
	}
	return &protoBuffers.User{
		Id:            u.ID.String(),
		Name:          u.Name,
		Position:      position,
		LastUpdated:   u.LastUpdated,
		CurrentSpidID: u.CurrentSpidID.String(),
	}, nil
}

func NewUser(name string, position gps.GlobalPosition) (*User, error) {
	if !position.IsValid() {
		return nil, fmt.Errorf("invalid position: %v", position)
	}
	return &User{
		ID:            uuid.New(),
		Name:          name,
		Position:      position,
		LastUpdated:   time.Now().Unix(),
		CurrentSpidID: uuid.Nil,
	}, nil
}

func (u User) Marshal() ([]byte, error) {
	return json.Marshal(u)
}

func MarshalUsers(users *Users) ([]byte, error) {
	return json.MarshalIndent(users, "", "    ")
}

func UnmarshalUsers(marshaledUsers []byte) (*Users, error) {
	var users Users
	err := json.Unmarshal(marshaledUsers, &users)
	return &users, err
}

func (u *User) UpdatePosition(position gps.GlobalPosition) error {
	if !position.IsValid() {
		return fmt.Errorf("invalid position: %v", position)
	}
	u.Position = position
	return nil
}
