package entities

import "../utils"

type Spid struct {
	ID string
	Location utils.GlobalPosition
	CurrentUserID string
	BatteryLevel float32
}
