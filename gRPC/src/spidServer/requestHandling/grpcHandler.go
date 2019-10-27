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
	h := Handler{
		DBManager: db.NewManager(basePath),
	}
	h.DBManager.LoadFromFile()
	return h
}
