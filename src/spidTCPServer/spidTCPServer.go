package main

import (
	"main/db"
	"main/entities"
)

const BasePath = "D:/GitReps/SD-College-Course-Work/src/spidTCPServer"

func main() {
	dbm := db.NewManager(BasePath)
	user := entities.NewUser("João")
	dbm.RegisterUser(user)

	//newUUID, err := uuid.Parse("5fc82534-f31f-4cdc-8814-3e1f47741395")
	//eh.HandleFatal(err)
	//fmt.Print(dbm.QueryUser(newUUID))

	//arguments := os.Args
	//if len(arguments) != 2 {
	//	_, filename, _, ok := runtime.Caller(1)
	//
	//	filenameSlice := strings.Split(filename, "/")
	//	filename = filenameSlice[len(filenameSlice)-1]
	//	if ok {
	//		fmt.Printf("%s usage: go run %s <port>\n", filename, filename)
	//	}
	//	return
	//}
	//
	//tcpServer.Listen(arguments[1])
}
