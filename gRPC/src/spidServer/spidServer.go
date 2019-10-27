package main

import (
	"fmt"
	"os"
	"runtime"
	"spidServer/grpcServer"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) != 4 {
		_, filename, _, ok := runtime.Caller(1)

		filenameSlice := strings.Split(filename, "/")
		filename = filenameSlice[len(filenameSlice)-1]
		if ok {
			fmt.Printf("%s usage: go run %s <port> <registrar address> <registrar port>\n", filename, filename)
		}
		return
	}
	port := arguments[1]
	registrarAddress := arguments[2]
	registrarPort := arguments[3]
	server := grpcServer.NewServer(port)
	server.Register(registrarAddress, registrarPort)
	server.Listen()
}
