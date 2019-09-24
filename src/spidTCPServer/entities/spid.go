package entities

import (
	"main/gps"
)

type Spid struct {
	ID            string `json:"id"`
	Location      gps.GlobalPosition `json:"location"`
	CurrentUserID string `json:"current_user_id"`
	BatteryLevel  float32 `json:"battery_level"`
}
