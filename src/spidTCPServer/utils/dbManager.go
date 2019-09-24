package utils

import (
	"io/ioutil"
	"log"
	"os"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

type DBManager struct {
	BasePath string
}

func (d DBManager) readFile(path string) []byte {
	path = d.BasePath + path
	file, err := os.Open(path)
	HandleFatal(err)
	defer HandleCloseFile(file, path)

	content, err := ioutil.ReadAll(file)
	HandleFatal(err)
	return content
}

func (d DBManager) writeToFile(path string, content []byte) {
	path = d.BasePath + path
	file, err := os.Open(path)
	HandleFatal(err)
	defer HandleCloseFile(file, path)

	_, err = file.Write(content)
	HandleFatal(err)
}

func (d DBManager) GetUsersFromFile() []byte {
	log.Print("Reading users.")
	return d.readFile(UsersDefaultLocation)
}

func (d DBManager) WriteUsersToFile(users []byte) {
	log.Printf("Writing users: %s", string(users))
	d.writeToFile(SpidsDefaultLocation, users)
}

func (d DBManager) GetSpidsFromFile() []byte {
	log.Print("Reading spids.")
	return d.readFile(SpidsDefaultLocation)
}

func (d DBManager) WriteSpidsToFile(spids []byte) {
	log.Printf("Writing users: %s", string(spids))
	d.writeToFile(SpidsDefaultLocation, spids)
}
