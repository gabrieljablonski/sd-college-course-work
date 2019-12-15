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
	if len(arguments) != 4 {
		filename := arguments[0]
		filename = strings.ReplaceAll(filename, "\\", "/")
		filename = path.Base(filename)
		fmt.Printf("%s usage: %s <port> <mapper address> <cluster endpoints>\n", filename, filename)
		return
	}
	port := arguments[1]
	mapperAddress := strings.Split(arguments[2], ":")

	mapperAddr := mapperAddress[0]
	mapperPort := mapperAddress[1]

	clusterEndpoints := strings.Split(arguments[3], ",")
	server := grpcServer.NewServer(clusterEndpoints, port)
	err = server.Handler.ConsensusManager.Recover()
	if err != nil {
		log.Fatalf("Failed to recover from consensus cluster: %v", err)
	}
	err = server.TryRegister(utils.IP{
		Address: mapperAddr,
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
