package requestHandling

import (
	"github.com/google/uuid"
	"log"
	"math"
	"spidServer/db"
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
	// number from 1 to `n` indicating position in global IP table
	Number		int
	ServerPoolSize int
	BaseDelta   int

	pb.SpidHandlerServer
	pb.UserHandlerServer
}

func NewHandler(basePath string) Handler {
	return Handler{DBManager: db.NewManager(basePath)}
}

func HostIsLocal(ip utils.IP) bool {
	return ip.Address == LocalHost
}

func (h *Handler) ClosestHost(targetServer int) utils.IP {
	if target, ok := h.IPMap[targetServer]; ok {
		return target
	}
	targetX := targetServer/h.BaseDelta
	targetY := targetServer % h.BaseDelta
	minDist := math.Inf(1)
	closestServer := utils.IP{}
	for n, ip := range h.IPMap {
		nx := n/h.BaseDelta
		ny := n % h.BaseDelta
		dist := math.Sqrt(math.Pow(float64(nx-targetX), 2) + math.Pow(float64(ny-targetY), 2))
		if dist <= minDist{
			minDist = dist
			closestServer = ip
		}
	}
	if closestServer.Address == "" {
		// this should never happen
		log.Fatalf("unexpected error finding closest host: targetServer=%s; minDist=%.2f; IPMap=%s",
					strconv.Itoa(targetServer), minDist, h.IPMap)
	}
	return closestServer
}

func (h *Handler) WhereIsPosition(position gps.GlobalPosition) utils.IP {
	h.BaseDelta = int(math.Round(math.Sqrt(float64(h.ServerPoolSize))))
	longitudeDelta := math.Ceil(360.0/float64(h.BaseDelta))
	latitudeDelta := math.Ceil(180.0/float64(h.BaseDelta))

	bLongitude := int(math.Floor(position.Longitude/longitudeDelta))
	bLatitude := int(math.Floor(position.Latitude/latitudeDelta))
	serverNumber := h.BaseDelta*bLatitude + bLongitude
	return h.ClosestHost(serverNumber)
}

func (h *Handler) WhereIsEntity(id uuid.UUID) utils.IP {
	// static rule for user and spid mapping
	//   -- all entities have a home server mapped by this rule
	//   -- all spids are also replicated to server
	//      mapped geographically so they can be easily found by users close by
	if len(h.IPMap) == 0 {
		log.Fatal("ip table is empty")
	}
	// uuid has uniform distribution
	serverNumber := id.ID() % uint32(len(h.IPMap))
	return h.ClosestHost(int(serverNumber))
}
