package utils

import (
	"log"
	"main/entities"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

type DBManager struct {
	FM FileManager
}

func NewManager(basePath string) DBManager {
	return DBManager{FM: FileManager{BasePath: basePath}}
}

func (d DBManager) GetUsersFromFile() []byte {
	log.Print("Reading users.")
	return d.FM.readFile(UsersDefaultLocation)
}

func (d DBManager) WriteUsersToFile(users []byte) {
	log.Printf("Writing users: %s", string(users))
	d.FM.writeToFile(SpidsDefaultLocation, users)
}

func (d DBManager) GetSpidsFromFile() []byte {
	log.Print("Reading spids.")
	return d.FM.readFile(SpidsDefaultLocation)
}

func (d DBManager) WriteSpidsToFile(spids []byte) {
	log.Printf("Writing users: %s", string(spids))
	d.FM.writeToFile(SpidsDefaultLocation, spids)
}

func (d DBManager) RegisterUser(user entities.User) {
	user.Register(d)
}