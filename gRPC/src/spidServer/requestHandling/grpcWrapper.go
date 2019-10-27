package requestHandling

import (
	"spidServer/db"
	pb "spidServer/requestHandling/protoBuffers"
)

type Wrapper struct {
	DBManager db.Manager
	pb.SpidHandlerServer
	pb.UserHandlerServer
}
