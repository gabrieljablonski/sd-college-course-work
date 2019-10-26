package db

import (
	"os"
	"spidServer/entities"
	"spidServer/utils"
)

const UsersDefaultLocation = "db" + string(os.PathSeparator) + "users.spd"
const SpidsDefaultLocation = "db" + string(os.PathSeparator) + "spids.spd"

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
