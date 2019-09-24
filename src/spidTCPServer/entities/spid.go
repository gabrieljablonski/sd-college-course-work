package entities

import (
	"spidTCPServer/gps"
)

type Spid struct {
	ID            string
	Location      gps.GlobalPosition
	CurrentUserID string
	BatteryLevel  float32
}
