package requestHandling

import (
	"github.com/google/uuid"
	"log"
	"math"
	"spidServer/db"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
)

const (
	LocalHost = "localhost"
)

type Handler struct {
	DBManager db.Manager
	IPTable	  []utils.IP
	// number from 1 to `n` indicating position in global IP table
	Number		int
	Boundaries  []gps.GlobalPosition

	pb.SpidHandlerServer
	pb.UserHandlerServer
}

func NewHandler(basePath string) Handler {
	return Handler{DBManager: db.NewManager(basePath)}
}

func HostIsLocal(ip utils.IP) bool {
	return ip.Address == LocalHost
}

func (h *Handler) CalculateBoundaries() {
	// visual representation of this algorithm applied for 9 divisions
	// can be seen in `global_boundaries_9_regions.png`.
	// the boundaries are represented by the red dots
	if len(h.IPTable) == 0 {
		log.Fatal("ip table is empty")
	}
	h.Boundaries = make([]gps.GlobalPosition, 0)
	// the result should be an integer
	// rounding to avoid situations like 1.99999... being truncated to 1
	baseDelta := int(math.Round(math.Sqrt(float64(len(s.Handler.IPTable)))))
	// using ceiling just to avoid rounding issues
	longitudeDelta := math.Ceil(360.0 / float64(baseDelta))
	latitudeDelta := math.Ceil(180.0 / float64(baseDelta))
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
			h.Boundaries = append(s.Boundaries, gps.GlobalPosition{
				Latitude:  lat,
				Longitude: lon,
			})
		}
	}
}

func (h *Handler) WhereIsPosition(position gps.GlobalPosition) utils.IP {
	if len(h.Boundaries) == 0 {
		log.Fatal("boundaries not calculated")
	}
	baseDelta := int(math.Round(math.Sqrt(float64(len(s.Boundaries)))))
	longitudeDelta := math.Ceil(360.0/float64(baseDelta))
	latitudeDelta := math.Ceil(180.0/float64(baseDelta))

	bLongitude := int(math.Floor(position.Longitude/longitudeDelta))
	bLatitude := int(math.Floor(position.Latitude/latitudeDelta))
	serverNumber := baseDelta*bLatitude + bLongitude
	return h.IPTable[serverNumber]
}

func (h *Handler) WhereIsEntity(id uuid.UUID) utils.IP {
	// static rule for user and spid mapping
	//   -- all entities have a home server mapped by this rule
	//   -- all spids are also replicated to server
	//      mapped geographically so they can be easily found by users close by
	if len(h.IPTable) == 0 {
		log.Fatal("ip table is empty")
	}
	// uuid has uniform distribution
	serverNumber := id.ID() % uint32(len(h.IPTable))
	return h.IPTable[serverNumber]
}
