package grpcServer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kardianos/osext"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
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
	DefaultPort     = "45678"
	GoogleDNS       = "8.8.8.8:80"

	TryUpdateIPMapPeriod = 1 * time.Second
)

type Server struct {
	ID         uuid.UUID
	Handler    requestHandling.Handler
	IP         utils.IP
	MapperIP   utils.IP
	Registered bool
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

func NewServer(clusterEndpoints []string, port string) Server {
	if port == "" {
		port = DefaultPort
	}
	basePath, err := osext.ExecutableFolder()
	errorHandling.HandleFatal(err)
	handler := requestHandling.NewHandler(clusterEndpoints, basePath)
	return Server{
		ID:      handler.ConsensusManager.DBManager.GetServerIDFromFile(),
		Handler: handler,
		IP: utils.IP{
			Address: GetOutboundIP(),
			Port:    port,
		},
	}
}

func (s *Server) TryRegister(mapperIP utils.IP) error {
	// placeholder solution
	// connect to server mapper to get a server number (from 0 to `n`-1)
	addr := mapperIP.String()
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		return fmt.Errorf("failed to register server at %s: %s", addr, err)
	}
	s.MapperIP = mapperIP
	s.Registered = true
	if s.ID == uuid.Nil {
		// this should happen only when the whole system is being setup
		log.Print("Nil server ID. Creating new ID.")
		s.ID = uuid.New()
		err = s.Handler.ConsensusManager.DBManager.WriteServerIDToFile(s.ID)
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
	serverPoolSize, err := strconv.Atoi(split[1])
	log.Printf("Server registered as %s.\n%d servers in total.", s.ID, serverPoolSize)
	s.Handler.ServerPoolSize = serverPoolSize
	s.Handler.BaseDelta = int(math.Round(math.Sqrt(float64(serverPoolSize))))
	return nil
}

func (s *Server) LoadIPMapFromFile() {
	ipMap, err := s.Handler.ConsensusManager.DBManager.GetIPMapFromFile()
	if err != nil {
		log.Fatalf("Failed to load IP map from file: %s", err)
	}
	s.Handler.IPMap = ipMap
}

func (s *Server) WaitRequestIPMapUpdate() {
	for {
		log.Print("Trying ip map update...")
		err := s.RequestIPMapUpdate()
		if err == nil {
			break
		}
		time.Sleep(TryUpdateIPMapPeriod)
	}
}

func (s *Server) RequestIPMapUpdate() error {
	if !s.Registered {
		log.Fatal("Server not registered")
	}
	addr := s.MapperIP.String()
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", addr, err)
	}
	request := fmt.Sprintf("REQUEST IP MAP %s\n", s.ID)
	log.Print("Sending request: ", request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Fatalf("Failed to send ip map request: %s", err)
	}
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to get response: %s", err)
	}
	log.Print("Received response: ", response)
	responseSplit := strings.Split(response, "\t")
	serverNumber, err := strconv.Atoi(responseSplit[0])
	if err != nil {
		log.Fatalf("Failed to parse server number `%s`: %s", serverNumber, err)
	}
	if serverNumber == -1 {
		msg := "ip map not ready"
		log.Print(msg)
		return fmt.Errorf(msg)
	}
	s.Handler.ServerNumber = serverNumber
	log.SetPrefix(fmt.Sprintf("#%d || %s >> ", serverNumber, s.IP.String()))
	responseIPMap := responseSplit[1]
	var ipMap map[int][]string
	err = json.Unmarshal([]byte(responseIPMap), &ipMap)
	if err != nil {
		log.Fatalf("Failed to parse ip map: %s", err)
	}
	//map size is dynamic now
	//if len(ipMap) != s.Handler.ServerPoolSize {
	//	msg := "unexpected ip map size, not ready"
	//	log.Print(msg)
	//	return fmt.Errorf(msg)
	//}
	s.Handler.IPMap = make(map[int][]utils.IP)
	for serverNumber, ips := range ipMap {
		if serverNumber == s.Handler.ServerNumber {
			s.Handler.IPMap[serverNumber] = []utils.IP{{
				Address: "localhost",
				Port:    s.IP.Port,
			}}
			continue
		}
		for _, ip := range ips {
			split := strings.Split(ip, ":")
			s.Handler.IPMap[serverNumber] = append(s.Handler.IPMap[serverNumber], utils.IP{
				Address: split[0],
				Port:    split[1],
			})
		}
	}
	log.Printf("Updated IP map: %s", s.Handler.IPMap)
	return s.Handler.ConsensusManager.DBManager.WriteIPMapToFile(s.Handler.IPMap)
}

func (s *Server) HandleRemoteEntities() {
	for len(s.Handler.IPMap) == 0 {
		// Wait until all servers are setup
		time.Sleep(time.Second)
	}
	log.Print("Handling remote entities...")
	s.handleRemoteUsers()
	s.handleRemoteSpids()
}

func (s *Server) handleRemoteUsers() {
	for _, user := range s.Handler.ConsensusManager.DBManager.Users.Users {
		err := s.Handler.HandleRemoteUser(user)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleRemoteSpids() {
	for _, spid := range s.Handler.ConsensusManager.DBManager.Spids.Spids {
		err := s.Handler.HandleRemoteSpid(spid)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) WatchChanges() {
	s.Handler.ConsensusManager.WatchChanges()
}

func (s *Server) Listen() {
	listener, err := net.Listen(DefaultProtocol, ":"+s.IP.Port)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err)
	}
	log.Printf("serving on %s:%s...", GetOutboundIP(), s.IP.Port)
	gs := grpc.NewServer()
	pb.RegisterSpidHandlerServer(gs, &s.Handler)
	// blocking call
	err = gs.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve connection: %s", err)
	}
}
