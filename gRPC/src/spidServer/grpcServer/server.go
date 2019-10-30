package grpcServer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"spidServer/errorHandling"
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
	"strconv"
	"strings"
)

const (
	DefaultProtocol = "tcp"
	DefaultPort = "43210"
)

type Server struct {
	ID 	        uuid.UUID
	Handler     requestHandling.Handler
	IP          utils.IP
	RegistrarIP utils.IP
	Registered  bool
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func NewServer(port string) Server {
	if port == "" {
		port = DefaultPort
	}
	basePath, err := os.Getwd()
	errorHandling.HandleFatal(err)
	handler := requestHandling.NewHandler(basePath)
	return Server{
		ID: handler.DBManager.GetServerID(),
		Handler: handler,
		IP: utils.IP{
			Address: GetOutboundIP(),
			Port:    port,
		},
	}
}

func (s *Server) Register(registrarIP utils.IP) {
	// placeholder solution
	// connect to server map registrar to get a server number (from 0 to `n`-1)
	addr := registrarIP.ToString()
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		log.Fatalf("failed to register server at %s: %s", addr, err)
	}
	s.RegistrarIP = registrarIP
	s.Registered = true
	if s.ID == uuid.Nil {
		// this should happen only when the whole system is being setup
		s.ID = uuid.New()
		err = s.Handler.DBManager.WriteServerID(s.ID)
		if err != nil {
			log.Fatalf("failed to save new server id: %s", err)
		}
	}
	request := fmt.Sprintf("REGISTER SERVER %s %s\n", s.ID.String(), s.IP.Port)
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("failed to send register request: %s", err)
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if response == "full\n" {
		log.Fatalf("all server slots are filled")
	}
	serverNumber, err := strconv.Atoi(strings.Trim(response, "\n"))
	if err != nil {
		log.Fatalf("failed to parse register response: %s", err)
	}
	log.Printf("Server registered as number %d", serverNumber)
	s.Handler.Number = serverNumber
}

func (s *Server) UpdateIPTable() error {
	if !s.Registered {
		log.Fatal("server not registered")
	}
	addr := s.RegistrarIP.ToString()
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		log.Fatalf("failed to connect to %s: %s", addr, err)
	}
	request := fmt.Sprintf("REQUEST IP TABLE\n")
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("failed to send ip table request: %s", err)
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("failed to get ip table response: %s", err)
	}
	var ipTableList []string
	err = json.Unmarshal([]byte(response), &ipTableList)
	if err != nil {
		log.Fatalf("failed to parse ip table: %s", err)
	}
	if len(ipTableList) == 0 {
		msg := "ip table not ready"
		log.Print(msg)
		return fmt.Errorf(msg)
	}
	s.Handler.IPTable = make([]utils.IP, 0)
	for i, ip := range ipTableList {
		if i == s.Handler.Number {
			s.Handler.IPTable = append(
				s.Handler.IPTable,
				s.IP,
			)
			continue
		}
		split := strings.Split(ip, ":")
		s.Handler.IPTable = append(
			s.Handler.IPTable,
			utils.IP{
				Address: split[0],
				Port:    split[1],
			},
		)
	}
	log.Printf("Updated IP table: %s", s.Handler.IPTable)
	s.Handler.CalculateBoundaries()
	return nil
}

func (s *Server) Listen() {
	listener, err := net.Listen(DefaultProtocol, ":" + s.IP.Port)
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}
	gs := grpc.NewServer()
	pb.RegisterSpidHandlerServer(gs, &s.Handler)
	log.Printf("serving on %s:%s...", GetOutboundIP(), s.IP.Port)
	// blocking call
	err = gs.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
}
