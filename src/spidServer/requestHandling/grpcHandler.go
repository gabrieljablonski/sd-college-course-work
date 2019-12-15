package requestHandling

import (
	"github.com/google/uuid"
	"log"
	"math"
	"spidServer/consensus"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
	"strconv"
	"time"
)

const (
	LocalHost = "localhost"
	DefaultContextTimeout = 60*time.Second
)

type Handler struct {
	ConsensusManager consensus.Manager
	IPMap     map[int][]utils.IP
	// number from 1 to `n` indicating index in global IP map
	ServerNumber   int
	ServerPoolSize int
	BaseDelta      int

	pb.SpidHandlerServer
}

func NewHandler(clusterEndpoints []string, basePath string) Handler {
	return Handler{
		ConsensusManager: consensus.NewManager(clusterEndpoints, basePath),
	}
}

func IsHostLocal(ip utils.IP) bool {
	return ip.Address == LocalHost
}

func (h *Handler) HandleRemoteUser(user *entities.User) error {
	targetServerNumber := h.WhereIsPosition(user.Position)
	if targetServerNumber == h.ServerNumber {
		return nil
	}
	log.Printf("Handling user: %s", user)
	pbUser, err := user.ToProtoBufferEntity()
	if err != nil {
		return err
	}
	return h.addRemoteUser(pbUser)
}

func (h *Handler) HandleRemoteSpid(spid *entities.Spid) error {
	targetServerNumber := h.WhereIsPosition(spid.Position)
	if targetServerNumber == h.ServerNumber {
		return nil
	}
	log.Printf("Handling spid: %s", spid)
	pbSpid, err := spid.ToProtoBufferEntity()
	if err != nil {
		return err
	}
	return h.addRemoteSpid(pbSpid)
}

func (h *Handler) getClosestHost(targetServer int) []utils.IP {
	if len(h.IPMap) == 0 {
		log.Fatal("IP map not setup.")
	}
	if target, ok := h.IPMap[targetServer]; ok {
		log.Printf("Target %d in ip map: %s.", targetServer, target)
		return target
	}
	targetX := targetServer/h.BaseDelta
	targetY := targetServer % h.BaseDelta
	minDist := math.Inf(1)
	number := -1
	var closestServers []utils.IP
	for n, ip := range h.IPMap {
		if IsHostLocal(ip[0]) {
			continue
		}
		nx := n/h.BaseDelta
		ny := n % h.BaseDelta
		dist := math.Sqrt(math.Pow(float64(nx-targetX), 2) + math.Pow(float64(ny-targetY), 2))
		if dist <= minDist {
			minDist = dist
			number = n
			closestServers = ip
		}
		if dist == 0 {
			break
		}
	}
	if len(closestServers) == 0 {
		// this should never happen
		log.Fatalf("Unexpected error finding closest host: targetServer=%s; minDist=%.2f; IPMap=%s",
					strconv.Itoa(targetServer), minDist, h.IPMap)
	}
	log.Printf("Target %d not in ip map, closest is %d: %s.", targetServer, number, closestServers)
	return closestServers
}

func (h *Handler) WhereIsPosition(position gps.GlobalPosition) int {
	log.Printf("Calculating where is %s.", position)
	longitudeDelta := math.Ceil(360.0/float64(h.BaseDelta))
	latitudeDelta := math.Ceil(180.0/float64(h.BaseDelta))
	position.Latitude += 90
	position.Longitude += 180
	bLongitude := int(math.Floor(position.Longitude/longitudeDelta))
	bLatitude := int(math.Floor(position.Latitude/latitudeDelta))
	if bLongitude >= h.BaseDelta {
		bLongitude = h.BaseDelta - 1
	}
	if bLatitude >= h.BaseDelta {
		bLatitude = h.BaseDelta - 1
	}
	serverNumber := h.BaseDelta*bLatitude + bLongitude
	log.Printf("Server number is %d.", serverNumber)
	return serverNumber
}

func (h *Handler) WhereIsEntity(id uuid.UUID) int {
	// static rule for user and spid mapping
	//   -- all entities have a home server mapped by this rule
	//   -- all spids are also replicated to server
	//      mapped geographically so they can be easily found by users close by
	log.Printf("Calculating where is %s (%d).", id, id.ID())
	if len(h.IPMap) == 0 {
		log.Fatal("ip map is empty")
	}
	// uuid has uniform distribution
	serverNumber := id.ID() % uint32(h.ServerPoolSize)
	log.Printf("Server number is %d.", serverNumber)
	return int(serverNumber)
}
