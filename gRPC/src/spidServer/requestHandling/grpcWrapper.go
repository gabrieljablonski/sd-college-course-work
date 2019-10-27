package requestHandling

import (
	"spidServer/db"
	pb "spidServer/requestHandling/spidProtoBuffers"
)

type Wrapper struct {
	DBManager db.Manager
	pb.SpidHandlerServer
	pb.UserHandlerServer
}
