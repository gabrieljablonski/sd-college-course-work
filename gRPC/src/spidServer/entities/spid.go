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

type LockInfo struct {
	Override bool `json:"override"`  // override user control
	Pending bool `json:"pending"`
	State string `json:"state"`
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

type Spid struct {
	ID            uuid.UUID          `json:"id"`
	BatteryLevel  uint32              `json:"battery_level"`
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

func (s Spid) ToProtoBufferEntity() *protoBuffers.SpidMinimal {
	return &protoBuffers.SpidMinimal{
		Id:           s.ID.String(),
		BatteryLevel: s.BatteryLevel,
		Location:     s.Location.ToProtoBufferEntity(),
		LockState:    s.Lock.State,
	}
}

func NewSpid() Spid {
	return Spid{
		ID: uuid.New(),
		BatteryLevel: 100,
		Lock: LockInfo{false, false, "locked"},
		Location: gps.NullPosition(),
		LastUpdated: time.Unix(0,0),
		CurrentUserID: uuid.Nil,
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

func (s *Spid) UpdateLockState(lockState string) error {
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
