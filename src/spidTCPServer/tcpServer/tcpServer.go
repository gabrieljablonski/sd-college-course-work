package tcpServer

import (
	"bufio"
	"log"
	eh "main/errorHandling"
	"net"
	"strings"
)

const DefaultEndConnection = "END CONNECTION"

func handleConnection(c net.Conn, incomingRequests, outgoingResponses chan string) {
	remoteAddr := c.RemoteAddr().String()
	log.Printf("Serving %s\n", remoteAddr)
	for {
		incomingData, err := bufio.NewReader(c).ReadString('\n')
		eh.HandleFatal(err)

		request := strings.TrimSpace(string(incomingData))
		if request == DefaultEndConnection {
			log.Printf("Ending connection with %s.\n", remoteAddr)
			break
		}
		incomingRequests <- request
		response := <-outgoingResponses
		_, err = c.Write([]byte(response + "\n"))
		eh.HandleFatal(err)
	}
	err := c.Close()
	eh.HandleFatal(err)
}

func Listen(port string, incomingRequests, outgoingResponses chan string) {
	port = ":" + port
	listener, err := net.Listen("tcp4", port)

	eh.HandleFatal(err)
	defer eh.HandleCloseListener(listener)

	log.Printf("Listing for connections on port %s...\n", port[1:])
	for {
		conn, err := listener.Accept()
		eh.HandleFatal(err)
		go handleConnection(conn, incomingRequests, outgoingResponses)
	}
}