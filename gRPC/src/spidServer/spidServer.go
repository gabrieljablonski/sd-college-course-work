package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"spidServer/requestHandling/grpcWrapper"
	pb "spidServer/requestHandling/grpcWrapper/spidProtoBuffers"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("failed", err)
	}
	s := grpc.NewServer()
	pb.RegisterSpidHandlerServer(s, &grpcWrapper.Wrapper{})
	log.Print("serving...")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal("failed", err)
	}
}