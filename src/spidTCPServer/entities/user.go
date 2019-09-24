package entities

import (
	"../utils"
)

type User struct {
	ID string
	Location utils.GlobalPosition
	CurrentSpidID string
}
