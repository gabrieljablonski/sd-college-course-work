package grpcWrapper

import (
	pb "spidServer/requestHandling/grpcWrapper/spidProtoBuffers"
)

type Wrapper struct {
	pb.SpidHandlerServer
	pb.UserHandlerServer
}
