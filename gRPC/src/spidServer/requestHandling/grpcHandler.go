package requestHandling

import (
	"spidServer/db"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
)

type Handler struct {
	DBManager db.Manager
	IPTable	  []utils.IP

	pb.SpidHandlerServer
	pb.UserHandlerServer
}

func NewHandler(basePath string) Handler {
	return Handler{DBManager: db.NewManager(basePath)}
}
