package db

import (
	"main/utils"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

type Manager struct {
	FM utils.FileManager
}

func NewManager(basePath string) Manager {
	return Manager{FM: utils.FileManager{BasePath: basePath}}
}
