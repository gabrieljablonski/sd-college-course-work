package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"spidServer/grpcServer"
	"spidServer/utils"
	"strings"
	"time"
)

const (
	UpdateIPMapPeriod = 1*time.Second
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
	server.Register(utils.IP{
		Address: registrarAddress,
		Port:    registrarPort,
	})
	for err := fmt.Errorf(""); err != nil; err = server.UpdateIPMap() {
		time.Sleep(UpdateIPMapPeriod)
		log.Print("Trying ip map update...")
	}
	server.Listen()
}
