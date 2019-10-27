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
