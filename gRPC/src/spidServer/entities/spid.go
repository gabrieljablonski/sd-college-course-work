package entities

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"main/gps"
	"time"
)

type LockInfo struct {
	Override bool `json:"override"`  // override user control
	Pending bool `json:"pending"`
	State string `json:"state"`
}

type Spid struct {
	ID            uuid.UUID          `json:"id"`
	BatteryLevel  uint8              `json:"battery_level"`
	Lock          LockInfo           `json:"lock"`
	Location      gps.GlobalPosition `json:"location"`
	LastUpdated   time.Time          `json:"last_updated"`
	CurrentUserID uuid.UUID          `json:"current_user_id"`
}

func (s Spid) ToString() string {
	ss, err := json.Marshal(s)
	if err != nil {
		log.Printf("Failed to convert spid to string: %s", err)
		return ""
	}
	return string(ss)
}

type Spids struct {
	Spids map[uuid.UUID]Spid `json:"spids"`
}

func (s Spids) ToString() string {
	ss, err := json.Marshal(s)
	if err != nil {
		log.Printf("Failed to convert spids to string: %s", err)
		return ""
	}
	return string(ss)
}

func NewSpid() Spid {
	return Spid{
		uuid.New(),
		100,
		LockInfo{false, false, "locked"},
		gps.NullPosition(),
		time.Unix(0,0),
		uuid.Nil,
	}
}

func IsValidLockState(lockState string) bool {
	switch lockState {
	case
		"locked",
		"unlocked":
		return true
	}
	return false
}

func (s Spid) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func MarshalSpids(spids Spids) ([]byte, error) {
	return json.MarshalIndent(spids, "", "    ")
}

func UnmarshalSpids(marshaledSpids []byte) (Spids, error) {
	var spids Spids
	err := json.Unmarshal(marshaledSpids, &spids)
	return spids, err
}

func (s *Spid) UpdateLocation(position gps.GlobalPosition) {
	s.Location = position
	s.LastUpdated = time.Now()
}

func (s *Spid) UpdateLockState(lockState string, userID uuid.UUID) error {
	if userID != s.CurrentUserID {
		return fmt.Errorf("user with id %s not associated with spid", userID)
	}
	if !IsValidLockState(lockState) {
		return fmt.Errorf("invalid lock state")
	}
	if s.Lock.Override {
		return fmt.Errorf("change lock state is unavailable")
	}
	s.Lock.State = lockState
	s.Lock.Pending = true
	return nil
}
