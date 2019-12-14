package main

import (
	"fmt"
	"github.com/kardianos/osext"
	"io"
	"log"
	"os"
	"path"
	"spidServer/db"
	"spidServer/errorHandling"
	"spidServer/grpcServer"
	"spidServer/utils"
	"strings"
)

func main() {
	basePath, err := osext.ExecutableFolder()
	errorHandling.HandleFatal(err)
	logPath := basePath + db.Sep + "logs.spd"
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wtr := io.MultiWriter(os.Stdout,f)
	log.SetFlags(log.LstdFlags)
	log.SetOutput(wtr)

	arguments := os.Args
	if len(arguments) != 5 {
		filename := arguments[0]
		filename = strings.ReplaceAll(filename, "\\", "/")
		filename = path.Base(filename)
		fmt.Printf("%s usage: %s <port> <mapper address> <mapper port> <cluster endpoints>\n", filename, filename)
		return
	}
	port := arguments[1]
	mapperAddress := arguments[2]
	mapperPort := arguments[3]
	clusterEndpoints := strings.Split(arguments[4], ",")
	server := grpcServer.NewServer(clusterEndpoints, port)
	err = server.TryRegister(utils.IP{
		Address: mapperAddress,
		Port:    mapperPort,
	})
	if err != nil {
		log.Printf("Failed to register to server mapper: %s", err)
		log.Print("Trying to load IP map from file")
		server.LoadIPMapFromFile()
	} else {
		server.WaitRequestIPMapUpdate()
	}
	go server.HandleRemoteEntities()
	go server.WatchChanges()
	server.Listen()
}
