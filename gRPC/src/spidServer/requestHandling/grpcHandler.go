package requestHandling

import (
	"github.com/google/uuid"
	"log"
	"math"
	"spidServer/db"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
	"strconv"
)

const (
	LocalHost = "localhost"
)

type Handler struct {
	DBManager db.Manager
	IPMap     map[int]utils.IP
	// number from 1 to `n` indicating index in global IP map
	Number		int
	ServerPoolSize int
	BaseDelta   int

	pb.SpidHandlerServer
	pb.UserHandlerServer
}

func NewHandler(basePath string) Handler {
	return Handler{DBManager: db.NewManager(basePath)}
}

func IsHostLocal(ip utils.IP) bool {
	return ip.Address == LocalHost
}

func (h *Handler) HandleRemoteUser(user *entities.User) error {
	ip := h.WhereIsPosition(user.Position)
	if IsHostLocal(ip) {
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
	ip := h.WhereIsPosition(spid.Position)
	if IsHostLocal(ip) {
		return nil
	}
	log.Printf("Handling spid: %s", spid)
	pbSpid, err := spid.ToProtoBufferEntity()
	if err != nil {
		return err
	}
	return h.addRemoteSpid(pbSpid)
}

func (h *Handler) getClosestHost(targetServer int) utils.IP {
	if target, ok := h.IPMap[targetServer]; ok {
		log.Printf("Target %d in ip map: %s.", targetServer, target)
		return target
	}
	targetX := targetServer/h.BaseDelta
	targetY := targetServer % h.BaseDelta
	minDist := math.Inf(1)
	number := -1
	closestServer := utils.IP{}
	for n, ip := range h.IPMap {
		nx := n/h.BaseDelta
		ny := n % h.BaseDelta
		dist := math.Sqrt(math.Pow(float64(nx-targetX), 2) + math.Pow(float64(ny-targetY), 2))
		if dist <= minDist {
			minDist = dist
			number = n
			closestServer = ip
		}
		if dist == 0 {
			break
		}
	}
	if closestServer.Address == "" {
		// this should never happen
		log.Fatalf("Unexpected error finding closest host: targetServer=%s; minDist=%.2f; IPMap=%s",
					strconv.Itoa(targetServer), minDist, h.IPMap)
	}
	log.Printf("Target %d not in ip map, closest is %d: %s.", targetServer, number, closestServer)
	return closestServer
}

func (h *Handler) WhereIsPosition(position gps.GlobalPosition) utils.IP {
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
	return h.getClosestHost(serverNumber)
}

func (h *Handler) WhereIsEntity(id uuid.UUID) utils.IP {
	// static rule for user and spid mapping
	//   -- all entities have a home server mapped by this rule
	//   -- all spids are also replicated to server
	//      mapped geographically so they can be easily found by users close by
	log.Printf("Calculating where is %s (%d).", id, id.ID())
	if len(h.IPMap) == 0 {
		log.Fatal("ip map is empty")
	}
	// uuid has uniform distribution
	serverNumber := id.ID() % uint32(len(h.IPMap))
	log.Printf("Server number is %d.", serverNumber)
	return h.getClosestHost(int(serverNumber))
}
