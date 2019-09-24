package tcpServer

import (
	"bufio"
	"log"
	"net"
	"spidTCPServer/utils"
	"strings"
)

const DefaultResponse = "CONNECTION SUCCESSFUL"
const DefaultEndConnection = "END CONNECTION"

func handleConnection(c net.Conn) {
	remoteAddr := c.RemoteAddr().String()
	log.Printf("Serving %s\n", remoteAddr)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		utils.HandleFatal(err)

		temp := strings.TrimSpace(string(netData))

		switch temp {
			case DefaultEndConnection:
				log.Printf("Ending connection with %s.\n", remoteAddr)
				break
		}
		if temp == DefaultEndConnection {
			break
		}

		result := DefaultResponse
		_, err = c.Write([]byte(string(result)+"\n"))
		utils.HandleFatal(err)
	}
	err := c.Close()
	utils.HandleFatal(err)
}

func Listen(port string) {
	port = ":" + port
	listener, err := net.Listen("tcp4", port)

	utils.HandleFatal(err)
	defer utils.HandleCloseListener(listener)

	log.Printf("Waiting for connection of port %s...\n", port[1:])
	for {
		conn, err := listener.Accept()
		utils.HandleFatal(err)
		go handleConnection(conn)
	}
}