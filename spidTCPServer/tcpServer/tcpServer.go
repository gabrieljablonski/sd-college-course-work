package tcpServer

import (
	"bufio"
	"crypto/md5"
	"log"
	eh "main/errorHandling"
	"main/requestHandling"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultEndConnection = "END CONNECTION"
	DefaultTimeout       = 2000*time.Millisecond
)

func handleConnection(c net.Conn, handler requestHandling.Handler) {
	remoteAddr := c.RemoteAddr().String()
	log.Printf("Serving %s.", remoteAddr)
	for {
		incomingData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Printf("Connection with %s lost.", remoteAddr)
			return
		}

		incomingMessage := strings.TrimSpace(string(incomingData))
		if incomingMessage == DefaultEndConnection {
			log.Printf("Ending connection requested by %s.", remoteAddr)
			_, _ = c.Write([]byte("ENDING CONNECTION"))
			break
		}
		requestMessage := requestHandling.GenericMessage{
			Message: incomingMessage,
			Received: time.Now(),
			Sum:     md5.Sum([]byte(incomingMessage+strconv.FormatInt(time.Now().UnixNano(), 10))),
		}
		log.Printf("%s: Logging request: `%s`", remoteAddr, requestMessage)
		err = handler.QueueRequest(requestMessage)
		if err != nil {
			log.Printf("Unexpected failure processing request. Ending connection with %s.\n`%s`", remoteAddr, requestMessage)
			break
		}
		log.Printf("Request logged.")
		log.Printf("%s: Waiting response...", remoteAddr)
		responseMessage, err := handler.GetResponse(requestMessage, DefaultTimeout)
		if err != nil {
			log.Printf("Error getting response: `%s`", err)
		}
		log.Printf("%s: Sending response.", remoteAddr)
		_, err = c.Write([]byte(responseMessage.Message + "\n"))
		if err != nil {
			log.Printf("Connection with %s lost.", remoteAddr)
			return
		}
	}
	err := c.Close()
	eh.HandleFatal(err)
}

func Listen(port string, handler requestHandling.Handler) {
	port = ":" + port
	listener, err := net.Listen("tcp4", port)

	eh.HandleFatal(err)
	defer eh.HandleCloseListener(listener)

	log.Printf("Listing for connections on port %s...\n", port[1:])
	for {
		conn, err := listener.Accept()
		eh.HandleFatal(err)
		go handleConnection(conn, handler)
	}
}