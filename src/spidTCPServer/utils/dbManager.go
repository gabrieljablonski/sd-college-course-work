package utils

import (
	"io/ioutil"
	"log"
	"os"
)

const UsersDefaultLocation = "/db/users.spd"
const SpidsDefaultLocation = "/db/spids.spd"

func readFile(path string) []byte {
	file, err := os.Open(path)
	HandleFatal(err)
	defer HandleCloseFile(file, path)

	content, err := ioutil.ReadAll(file)
	HandleFatal(err)
	return content
}

func writeToFile(path string, content []byte) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	HandleFatal(err)
	defer HandleCloseFile(file, path)

	_, err = file.Write(content)
	HandleFatal(err)
}

func GetUsersFromFile() []byte {
	log.Print("Reading users.")
	return readFile(UsersDefaultLocation)
}

func WriteToUsersFile(users []byte) {
	log.Printf("Writing users: %s", string(users))
	writeToFile(SpidsDefaultLocation, users)
}

func GetSpidsFromFile() []byte {
	log.Print("Reading spids.")
	return readFile(SpidsDefaultLocation)
}

func WriteToSpidsFile(spids []byte) bool {
	log.Printf("Writing users: %s", string(spids))
	writeToFile(SpidsDefaultLocation, spids)
}
