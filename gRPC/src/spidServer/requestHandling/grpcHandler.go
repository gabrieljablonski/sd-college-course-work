package requestHandling

import (
	"spidServer/db"
	"spidServer/grpcServer"
	pb "spidServer/requestHandling/protoBuffers"
)

type Handler struct {
	DBManager db.Manager
	IPTable	              []grpcServer.IP

	pb.SpidHandlerServer
	pb.UserHandlerServer
}

func NewHandler(basePath string) Handler {
	return Handler{DBManager: db.NewManager(basePath)}
}
