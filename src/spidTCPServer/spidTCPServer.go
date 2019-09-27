package main

import (
	"fmt"
	"main/errorHandling"
	"main/requestHandling"
	"main/tcpServer"
	"os"
	"runtime"
	"strings"
)

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
	basePath, err := os.Getwd()
	errorHandling.HandleFatal(err)
	handler := requestHandling.NewHandler(basePath)
	port := arguments[1]
	tcpServer.Listen(port, handler)
}
