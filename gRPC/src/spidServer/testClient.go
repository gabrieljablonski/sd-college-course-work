package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "spidServer/requestHandling/protoBuffers"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed connecting", err)
	}
	defer conn.Close()
	c := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RegisterSpid(ctx, &pb.RegisterSpidRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r)
}
