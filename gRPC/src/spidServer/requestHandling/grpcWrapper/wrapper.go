package grpcWrapper

import (
	"spidServer/requestHandling/grpc/spidPB"
	"spidServer/requestHandling/grpc/userPB"
)

type Wrapper struct {
	spidPB.SpidHandlerServer
	userPB.UserHandlerServer
}
