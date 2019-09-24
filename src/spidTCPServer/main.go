package spidTCPServer

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"./tcpServer"
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

	tcpServer.Listen(arguments[1])
}
