package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"time"
)

func main() {
	nServers := 9
	server := rand.Intn(nServers)
	host := fmt.Sprintf("localhost:%d", server+45678)
	log.Print("Host: ", host)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed connecting", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	c := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	spids := make([]*pb.Spid, 0)
	for i:=0; i<100; i++ {
		log.Print(i)
		spid, _ := entities.NewSpid(1, gps.GlobalPosition{
			Latitude:  -90 + 180*rand.Float64(),
			Longitude: -180 + 360*rand.Float64(),
		})
		pbSpid, _ := spid.ToProtoBufferEntity()
		log.Print(pbSpid)
		register, err := c.RegisterSpid(ctx, &pb.RegisterSpidRequest{
			Spid: pbSpid,
		})
		if err != nil {
			log.Print(i)
			log.Fatalf("Could not register: %s", err)
		}
		log.Printf("%s", register)
		spids = append(spids, register.Spid)
		getInfo, err := c.GetSpidInfo(ctx, &pb.GetSpidRequest{
			SpidID: register.Spid.Id,
		})
		if err != nil {
			log.Fatalf("Could not get info: %s", err)
		}
		log.Printf("%s", getInfo)
	}
	var a int
	_, _ = fmt.Scan(&a)
	for i, spid := range spids{
		for {
			log.Print(i)
			spid.Position = &pb.GlobalPosition{
				Latitude:  -90,
				Longitude: -180,
			}
			update, err := c.UpdateSpid(ctx, &pb.UpdateSpidRequest{
				Spid: spid,
			})
			if err != nil {
				log.Print(err)
			} else {
				log.Print(update)
				log.Print(update.Spid)
				break
			}
		}
	}
	//log.Printf("%s", getInfo)
	//getInfo.Spid.Position = &pb.GlobalPosition{
	//	Latitude:             90,
	//	Longitude:            180,
	//}
	//update, err := c.UpdateSpid(ctx, &pb.UpdateSpidRequest{
	//	Spid: getInfo.Spid,
	//})
	//if err != nil {
	//	log.Fatalf("Could not update spid: %s", err)
	//}
	//log.Printf("%s", update)
	//spids, err := c.GetRemoteSpids(ctx, &pb.GetRemoteSpidsRequest{
	//	Position: &pb.GlobalPosition{
	//		Latitude:  0,
	//		Longitude: 0,
	//	},
	//})
	//log.Print(spids)
}
