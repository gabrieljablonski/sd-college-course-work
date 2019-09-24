package utils

import (
	"encoding/json"
	"log"
	"os"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

func checkExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func RegisterUser(user interface{}) bool {
	u, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error creating user. %s", err)
		return false
	}
	log.Printf("Creating user: %s", u)

	if !checkExists(UsersDefaultLocation) {
		log.Printf("%s does not exist, creating file...", UsersDefaultLocation)
		newFile, err := os.Create(UsersDefaultLocation)
		HandleFatal(err)

		defer HandleCloseFile(newFile, UsersDefaultLocation)

	}
	return true
}

