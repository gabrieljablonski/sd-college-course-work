package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"time"
)

func newConnection(host string) (*grpc.ClientConn, context.Context, context.CancelFunc) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed connecting", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	return conn, ctx, cancel
}

func main() {
	host := fmt.Sprintf("localhost:%d", 8888)
	log.Print("Host: ", host)
	//spids := make([]*pb.Spid, 0)
	//for i:=0; i<100; i++ {
	//	go func(i int) {
	//		time.Sleep(time.Duration(rand.Float64())*time.Second)
	//		log.Print(i)
	//		conn, ctx, cancel := newConnection(host)
	//		defer func() {
	//			_ = conn.Close()
	//		}()
	//		defer cancel()
	//		c := pb.NewSpidHandlerClient(conn)
	//		spid, _ := entities.NewSpid(1, gps.Random())
	//		pbSpid, _ := spid.ToProtoBufferEntity()
	//		log.Print(pbSpid)
	//		register, err := c.RegisterSpid(ctx, &pb.RegisterSpidRequest{
	//			Spid: pbSpid,
	//		})
	//		if err != nil {
	//			log.Print(i)
	//			log.Fatalf("Could not register: %s", err)
	//		}
	//		log.Printf("%s", register)
	//		spids = append(spids, register.Spid)
	//		getInfo, err := c.GetSpidInfo(ctx, &pb.GetSpidRequest{
	//			SpidID: register.Spid.Id,
	//		})
	//		if err != nil {
	//			log.Fatalf("Could not get info: %s", err)
	//		}
	//		log.Printf("%s", getInfo)
	//	}(i)
	//}
	//var a int
	//_, _ = fmt.Scan(&a)
	conn, ctx, cancel := newConnection(host)
	defer func() {
		_ = conn.Close()
	}()
	defer cancel()
	u := pb.NewSpidHandlerClient(conn)
	user, _ := entities.NewUser("João", gps.Random())
	pbUser, _ := user.ToProtoBufferEntity()
	newUser, err := u.RegisterUser(ctx, &pb.RegisterUserRequest{
		User: pbUser,
	})
	if err != nil {
		log.Fatalf("%s", err)
	}
	remoteSpids, err := u.GetRemoteSpids(ctx, &pb.GetRemoteSpidsRequest{
		Position: newUser.User.Position,
	})
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Print(remoteSpids.MarshaledSpids)
	var spidss entities.Spids
	err = json.Unmarshal([]byte(remoteSpids.MarshaledSpids), &spidss)
	if err != nil {
		log.Fatalf("%s", err)
	}
	var spidID string
	for _, v := range spidss.Spids {
		spidID = v.ID.String()
	}
	associate, err := u.RequestAssociation(ctx, &pb.RequestAssociationRequest{
		UserID: pbUser.Id,
		SpidID: spidID,
	})
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Print(associate)
	associate, err = u.RequestAssociation(ctx, &pb.RequestAssociationRequest{
		UserID: pbUser.Id,
		SpidID: spidID,
	})
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Print(associate)
	//var a int
	//_, _ = fmt.Scan(&a)
	//for i, spid := range spids{
	//	for {
	//		log.Print(i)
	//		spid.Position = &pb.GlobalPosition{
	//			Latitude:  -90,
	//			Longitude: -180,
	//		}
	//		update, err := c.UpdateSpid(ctx, &pb.UpdateSpidRequest{
	//			Spid: spid,
	//		})
	//		if err != nil {
	//			log.Print(err)
	//		} else {
	//			log.Print(update)
	//			log.Print(update.Spid)
	//			break
	//		}
	//	}
	//}
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
