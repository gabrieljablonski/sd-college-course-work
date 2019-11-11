package main

import (
	"fmt"
	"github.com/kardianos/osext"
	"io"
	"log"
	"os"
	"runtime"
	"spidServer/db"
	"spidServer/errorHandling"
	"spidServer/grpcServer"
	"spidServer/utils"
	"strings"
)

const (
)

func main() {
	basePath, err := osext.ExecutableFolder()
	errorHandling.HandleFatal(err)
	path := basePath + db.Sep + "logs.spd"
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wtr := io.MultiWriter(os.Stdout,f)
	log.SetOutput(wtr)

	arguments := os.Args
	if len(arguments) != 4 {
		_, filename, _, ok := runtime.Caller(1)

		filenameSlice := strings.Split(filename, "/")
		filename = filenameSlice[len(filenameSlice)-1]
		if ok {
			fmt.Printf("%s usage: go run %s <port> <mapper address> <mapper port>\n", filename, filename)
		}
		return
	}
	port := arguments[1]
	mapperAddress := arguments[2]
	mapperPort := arguments[3]
	server := grpcServer.NewServer(port)
	err = server.TryRegister(utils.IP{
		Address: mapperAddress,
		Port:    mapperPort,
	})
	if err != nil {
		server.LoadIPMapFromFile()
	} else {
		server.WaitRequestIPMapUpdate()
	}
	go server.HandleRemoteEntities()
	server.Listen()
}
