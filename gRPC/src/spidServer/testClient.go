package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"spidServer/requestHandling"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed connecting", err)
	}
	defer conn.Close()
	c := requestHandling.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RegisterSpid(ctx, &requestHandling.ClientRequest{
		Id:   "1111",
		Type: requestHandling.RegisterSpid,
		Body: "{\"message\": \"gimme a new spid!\"}",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r)
}
