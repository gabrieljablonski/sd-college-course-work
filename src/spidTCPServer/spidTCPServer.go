package main

import (
	"fmt"
	"main/requestHandling"
	"os"
	"runtime"
	"strings"
)

const BasePath = "D:/GitReps/SD-College-Course-Work/src/spidTCPServer"

func main() {
	arguments := os.Args
	if len(arguments) != 2 {
		_, filename, _, ok := runtime.Caller(1)

		filenameSlice := strings.Split(filename, "/")
		filename = filenameSlice[len(filenameSlice)-1]
		if ok {
			fmt.Printf("%s usage: go run %s <port>\n", filename, filename)
		}
		return
	}
	handler := requestHandling.NewHandler(BasePath)

	fmt.Print(handler.ProcessRequest(
		`{
			"id": "88b7dd20-a91c-4693-80b8-a7683ce3ee97",
			"type": "REGISTER USER",
			"body": {
				"name": "João da Silva"
			}
		}`))
	//port := arguments[1]
	//handler.Listen(port)
}
