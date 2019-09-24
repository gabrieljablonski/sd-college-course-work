package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

const DefaultResponse = "CONNECTION SUCCESSFUL"
const DefaultEndConnection = "END CONNECTION"

func handleError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if !handleError(err) {
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == DefaultEndConnection {
			break
		}

		result := DefaultResponse
		_, err = c.Write([]byte(string(result)+"\n"))
		handleError(err)
	}
	err := c.Close()
	handleError(err)
}

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

	port := ":" + arguments[1]
	listener, err := net.Listen("tcp4", port)

	if !handleError(err) {
		return
	}

	defer listener.Close()

	fmt.Printf("Waiting for connection of port %s...\n", port[1:])
	for {
		conn, err := listener.Accept()
		if !handleError(err) {
			return
		}
		go handleConnection(conn)
	}
}
