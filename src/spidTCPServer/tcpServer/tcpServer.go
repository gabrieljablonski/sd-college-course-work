package tcpServer

import (
	"bufio"
	"log"
	eh "main/errorHandling"
	"net"
	"strings"
)

const DefaultResponse = "CONNECTION SUCCESSFUL"
const DefaultEndConnection = "END CONNECTION"

func handleConnection(c net.Conn) {
	remoteAddr := c.RemoteAddr().String()
	log.Printf("Serving %s\n", remoteAddr)
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		eh.HandleFatal(err)

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
		eh.HandleFatal(err)
	}
	err := c.Close()
	eh.HandleFatal(err)
}

func Listen(port string) {
	port = ":" + port
	listener, err := net.Listen("tcp4", port)

	eh.HandleFatal(err)
	defer eh.HandleCloseListener(listener)

	log.Printf("Listing for connections on port %s...\n", port[1:])
	for {
		conn, err := listener.Accept()
		eh.HandleFatal(err)
		go handleConnection(conn)
	}
}