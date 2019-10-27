package requestHandling

import (
	"spidServer/db"
	pb "spidServer/requestHandling/protoBuffers"
)

type Handler struct {
	DBManager db.Manager
	pb.SpidHandlerServer
	pb.UserHandlerServer
}

func NewHandler(basePath string) Handler {
	return Handler{DBManager: db.NewManager(basePath)}
}
