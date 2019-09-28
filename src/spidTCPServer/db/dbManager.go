package db

import (
	"main/entities"
	"main/utils"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

type Manager struct {
	FileManager utils.FileManager
	Users       entities.Users
	Spids       entities.Spids
}

func NewManager(basePath string) Manager {
	return Manager{FileManager: utils.FileManager{BasePath: basePath}}
}

func (m *Manager) LoadFromFile() {
	m.Users = m.GetUsersFromFile()
	m.Spids = m.GetSpidsFromFile()
}
