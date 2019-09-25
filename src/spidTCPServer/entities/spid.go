package entities

import (
	"encoding/json"
	"github.com/google/uuid"
	eh "main/errorHandling"
	"main/gps"
	"time"
)

type Spid struct {
	ID            uuid.UUID          `json:"id"`
	BatteryLevel  uint8              `json:"battery_level"`
	LockState     string             `json:"lock_state"`
	Location      gps.GlobalPosition `json:"location"`
	LastUpdated   time.Time          `json:"last_updated"`
	CurrentUserID string             `json:"current_user_id"`
}

type Spids struct {
	Spids map[uuid.UUID]Spid `json:"spids"`
}

func NewSpid() Spid {
	return Spid{
		uuid.New(),
		100,
		"locked",
		gps.NullPosition(),
		time.Unix(0,0),
		"",
	}
}

func (s Spid) Marshal() []byte {
	marshaledSpid, err := json.Marshal(s)
	eh.HandleFatal(err)
	return marshaledSpid
}

func MarshalSpids(spids Spids) []byte {
	marshaledSpids, err := json.MarshalIndent(spids, "", "    ")
	eh.HandleFatal(err)
	return marshaledSpids
}

func UnmarshalSpids(marshaledSpids []byte) Spids {
	var spids Spids
	err := json.Unmarshal(marshaledSpids, &spids)
	eh.HandleFatal(err)
	return spids
}

func (s Spid) UpdateLocation(position gps.GlobalPosition) {
	s.Location = position
	s.LastUpdated = time.Now()
}
