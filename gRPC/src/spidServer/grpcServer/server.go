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
	"time"
)

const (
	DefaultProtocol = "tcp"
	DefaultPort = "45678"

	GoogleDNS = "8.8.8.8:80"
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
	conn, err := net.Dial("udp", GoogleDNS)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()
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
	addr := registrarIP.String()
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		log.Fatalf("Failed to register server at %s: %s", addr, err)
	}
	s.RegistrarIP = registrarIP
	s.Registered = true
	if s.ID == uuid.Nil {
		// this should happen only when the whole system is being setup
		log.Print("Nil server ID. Creating new ID.")
		s.ID = uuid.New()
		err = s.Handler.DBManager.WriteServerID(s.ID)
		if err != nil {
			log.Fatalf("Failed to save new server id: %s", err)
		}
	}
	request := fmt.Sprintf("REGISTER SERVER %s %s\n", s.ID, s.IP.Port)
	log.Print("Sending request: ", request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("Failed to send register request: %s", err)
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if response == "full\n" {
		log.Fatalf("All server slots are filled")
	}
	log.Print("Received response: ", response)
	split := strings.Split(strings.Trim(response, "\n"), " ")
	serverNumber, err := strconv.Atoi(split[0])
	if err != nil {
		log.Fatalf("Failed to parse register response: %s", err)
	}
	serverPoolSize, err := strconv.Atoi(split[1])
	log.Printf("Server registered as number %d.\n%d servers in total.", serverNumber, serverPoolSize)
	s.Handler.Number = serverNumber
	s.Handler.ServerPoolSize = serverPoolSize
}

func (s *Server) UpdateIPMap() error {
	if !s.Registered {
		log.Fatal("Server not registered")
	}
	addr := s.RegistrarIP.String()
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", addr, err)
	}
	request := fmt.Sprintf("REQUEST IP MAP %d\n", s.Handler.Number)
	log.Print("Sending request: ", request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("Failed to send ip map request: %s", err)
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to get ip map response: %s", err)
	}
	log.Print("Received response: ", response)
	var ipMap map[int]string
	err = json.Unmarshal([]byte(response), &ipMap)
	if err != nil {
		log.Fatalf("Failed to parse ip map: %s", err)
	}
	if len(ipMap) == 0 {
		msg := "ip map not ready"
		log.Print(msg)
		return fmt.Errorf(msg)
	}
	if len(ipMap) != s.Handler.ServerPoolSize {
		msg := "unexpected ip map size not ready"
		log.Print(msg)
		return fmt.Errorf(msg)
	}
	s.Handler.IPMap = make(map[int]utils.IP)
	for serverNumber, ip := range ipMap {
		if serverNumber == s.Handler.Number {
			s.Handler.IPMap[serverNumber] = utils.IP{
				Address: "localhost",
				Port:    s.IP.Port,
			}
			continue
		}
		split := strings.Split(ip, ":")
		s.Handler.IPMap[serverNumber] = utils.IP{
			Address: split[0],
			Port:    split[1],
		}
	}
	log.Printf("Updated IP map: %s", s.Handler.IPMap)
	return nil
}

func (s *Server) HandleRemoteEntities() {
	for ; len(s.Handler.IPMap) == 0; {
		// Wait until all servers are setup
		time.Sleep(time.Second)
	}
	log.Print("Handling remote entities...")
	s.handleRemoteUsers()
	s.handleRemoteSpids()
}

func (s *Server) handleRemoteUsers() {
	for _, user := range s.Handler.DBManager.Users.Users {
		_ = s.Handler.HandleRemoteUser(user)
	}
}

func (s *Server) handleRemoteSpids() {
	for _, spid := range s.Handler.DBManager.Spids.Spids {
		_ = s.Handler.HandleRemoteSpid(spid)
	}
}

func (s *Server) Listen() {
	listener, err := net.Listen(DefaultProtocol, ":" + s.IP.Port)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err)
	}
	gs := grpc.NewServer()
	pb.RegisterSpidHandlerServer(gs, &s.Handler)
	log.Printf("serving on %s:%s...", GetOutboundIP(), s.IP.Port)
	// blocking call
	err = gs.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve connection: %s", err)
	}
}
