package requestHandling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"time"
)

func (h *Handler) querySpid(spidID string) (*pb.Spid, error) {
	id, err := uuid.Parse(spidID)
	if err != nil {
		return nil, fmt.Errorf("invalid spid id: %s", err)
	}
	ip := h.WhereIsEntity(id)
	if HostIsLocal(ip) {
		spid, err := h.DBManager.QuerySpid(id)
		if err != nil {
			return nil, err
		}
		return spid.ToProtoBufferEntity(), err
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.GetSpidRequest{
		SpidID: id.String(),
	}
	response, err := client.GetSpidInfo(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Spid, err
}

func (h *Handler) registerSpid(batteryLevel uint32, location gps.GlobalPosition) (*pb.Spid, error) {
	spid := entities.NewSpid(batteryLevel, location)
	ip := h.WhereIsEntity(spid.ID)
	if HostIsLocal(ip) {
		err := h.DBManager.RegisterSpid(spid)
		if err != nil {
			return nil, err
		}
		return spid.ToProtoBufferEntity(), nil
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.RegisterSpidRequest{
		BatteryLevel: batteryLevel,
		Location: location.ToProtoBufferEntity(),
	}
	response, err := client.RegisterSpid(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Spid, nil
}

func (h *Handler) updateSpid(pbSpid *pb.Spid) error {
	id, err := uuid.Parse(pbSpid.Id)
	if err != nil {
		return err
	}
	ip := h.WhereIsEntity(id)
	if HostIsLocal(ip) {
		spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
		if err != nil {
			return err
		}
		return h.DBManager.UpdateSpid(spid)
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.UpdateSpidRequest{
		Spid: pbSpid,
	}
	_, err = client.UpdateSpid(ctx, request)
	return err
}

func (h *Handler) deleteSpid(pbSpid *pb.Spid) error {
	id, err := uuid.Parse(pbSpid.Id)
	if err != nil {
		return err
	}
	ip := h.WhereIsEntity(id)
	if HostIsLocal(ip) {
		spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
		if err != nil {
			return err
		}
		return h.DBManager.UpdateSpid(spid)
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.DeleteSpidRequest{
		SpidID: pbSpid.Id,
	}
	_, err = client.DeleteSpid(ctx, request)
	return err
}
