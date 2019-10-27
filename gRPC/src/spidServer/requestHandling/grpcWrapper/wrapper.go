package grpcWrapper

import (
	"spidServer/requestHandling"
	pb "spidServer/requestHandling/grpcWrapper/spidProtoBuffers"
)

type Wrapper struct {
	Handler requestHandling.Handler
	pb.SpidHandlerServer
	pb.UserHandlerServer
}
