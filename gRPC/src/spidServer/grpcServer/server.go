package grpcServer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"os"
	"spidServer/errorHandling"
	"spidServer/gps"
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
	"strconv"
	"strings"
)

const (
	DefaultProtocol = "tcp"
	DefaultPort = "5000"
)

type Server struct {
	ID 				    uuid.UUID
	Handler             requestHandling.Handler
	IP                  utils.IP
	// number from 1 to `n` indicating position in global IP table
	Number				int
	RegistrarConnection net.Conn
	Boundaries          []gps.GlobalPosition
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return string(localAddr.IP)
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

func (s *Server) CalculateBoundaries() {
	// visual representation of this algorithm applied for 9 divisions
	// can be seen in `global_boundaries_9_regions.png`.
	// the boundaries are represented by the red dots
	if len(s.Handler.IPTable) == 0 {
		log.Fatal("ip table is empty")
	}
	s.Boundaries = make([]gps.GlobalPosition, 0)
	// the result should be an integer
	// rounding to avoid situations like 1.99999... being truncated to 1
	baseDelta := int(math.Round(math.Sqrt(float64(len(s.Handler.IPTable)))))
	// using ceiling just to avoid rounding issues
	longitudeDelta := math.Ceil(360.0/float64(baseDelta))
	latitudeDelta := math.Ceil(180.0/float64(baseDelta))
	for i := 1; i <= baseDelta; i++ {
		for j := 1; j <= baseDelta; j++ {
			lat := -90.0 + float64(i)*latitudeDelta
			lon := -180.0 + float64(j)*longitudeDelta
			// fixing upper limits manually to avoid rounding issues
			if i == baseDelta {
				lat = 90.0
			}
			if j == baseDelta {
				lon = 180.0
			}
			s.Boundaries = append(s.Boundaries, gps.GlobalPosition{
				Latitude:  lat,
				Longitude: lon,
			})
		}
	}
}

func (s *Server) WhereIsPosition(position gps.GlobalPosition) utils.IP {
	if len(s.Boundaries) == 0 {
		log.Fatal("boundaries not calculated")
	}
	baseDelta := int(math.Round(math.Sqrt(float64(len(s.Boundaries)))))
	longitudeDelta := math.Ceil(360.0/float64(baseDelta))
	latitudeDelta := math.Ceil(180.0/float64(baseDelta))

	bLongitude := int(math.Floor(position.Longitude/longitudeDelta))
	bLatitude := int(math.Floor(position.Latitude/latitudeDelta))
	serverNumber := baseDelta*bLatitude + bLongitude
	return s.Handler.IPTable[serverNumber]
}

func (s *Server) WhereIsEntity(id uuid.UUID) utils.IP {
	// static rule for user and spid mapping
	//   -- all entities have a home server mapped by this rule
	//   -- all spids are also replicated to server
	//      mapped geographically so they can be easily found by users close by
	if len(s.Handler.IPTable) == 0 {
		log.Fatal("ip table is empty")
	}
	// uuid has uniform distribution
	serverNumber := id.ID() % uint32(len(s.Handler.IPTable))
	return s.Handler.IPTable[serverNumber]
}

func (s *Server) Register(registrarAddress, registrarPort string) {
	// placeholder solution
	// connect to server map registrar to get a server number (from 0 to `n`-1)
	addr := fmt.Sprintf("%s:%s", registrarAddress, registrarPort)
	conn, err := net.Dial(DefaultProtocol, addr)
	if err != nil {
		log.Fatalf("failed to register server at %s", addr)
	}
	s.RegistrarConnection = conn
	if s.ID == uuid.Nil {
		// this should happen only when the whole system is being setup
		s.ID = uuid.New()
		err = s.Handler.DBManager.WriteServerID(s.ID)
		if err != nil {
			log.Fatalf("failed to save new server id: %s", err)
		}
	}
	request := fmt.Sprintf("REGISTER SERVER %s %s\n", s.IP.Port, s.ID.String())
	_, err = s.RegistrarConnection.Write([]byte(request))
	if err != nil {
		log.Fatalf("failed to send register request: %s", err)
	}
	response, err := bufio.NewReader(s.RegistrarConnection).ReadString('\n')
	serverNumber, err := strconv.Atoi(response)
	if err != nil {
		log.Fatalf("failed to parse register response: %s", err)
	}
	s.Number = serverNumber
}

func (s *Server) UpdateIPTable() {
	if s.RegistrarConnection == nil {
		log.Fatal("server not registered")
	}
	request := fmt.Sprintf("REQUEST IP TABLE\n")
	_, err := s.RegistrarConnection.Write([]byte(request))
	if err != nil {
		log.Fatalf("failed to send ip table request: %s", err)
	}
	response, err := bufio.NewReader(s.RegistrarConnection).ReadString('\n')
	if err != nil {
		log.Fatalf("failed to get ip table response: %s", err)
	}
	var ipTableList []string
	err = json.Unmarshal([]byte(response), &ipTableList)
	if err != nil {
		log.Fatalf("failed to parse ip table: %s", err)
	}
	for _, ip := range ipTableList {
		split := strings.Split(ip, ":")
		s.Handler.IPTable = append(
			s.Handler.IPTable,
			utils.IP{
				Address: split[0],
				Port:    split[1],
			},
		)
	}
	s.GlobalDivisions = len(s.Handler.IPTable)
	s.ComputeBoundaries()
}

func (s *Server) Listen() {
	listener, err := net.Listen(DefaultProtocol, ":" + s.IP.Port)
	if err != nil {
		log.Fatalf("failed to create listener: %s", err)
	}
	gs := grpc.NewServer()
	pb.RegisterSpidHandlerServer(gs, &s.Handler)
	log.Printf("serving on %s%s...", GetOutboundIP(), s.IP.Port)
	// blocking call
	err = gs.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
}
