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

func (h *Handler) registerSpid(batteryLevel uint32, position gps.GlobalPosition) (*pb.Spid, error) {
	spid := entities.NewSpid(batteryLevel, position)
	ip := h.WhereIsEntity(spid.ID)
	if HostIsLocal(ip) {
		err := h.DBManager.RegisterSpid(spid)
		if err != nil {
			return nil, err
		}
		err = h.addRemoteSpid(spid.ToProtoBufferEntity())
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
	request := &pb.RegisterSpidRequest{
		BatteryLevel: batteryLevel,
		Position:     position.ToProtoBufferEntity(),
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
		err = h.DBManager.UpdateSpid(spid)
		if err != nil {
			return err
		}
		return h.updateRemoteSpid(pbSpid)
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
		err = h.DBManager.DeleteSpid(spid)
		if err != nil {
			return err
		}
		return h.removeRemoteSpid(pbSpid)
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

func (h *Handler) addRemoteSpid(pbSpid *pb.Spid) error {
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbSpid.Position))
	if HostIsLocal(ip) {
		spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
		if err != nil {
			return err
		}
		return h.DBManager.AddRemoteSpid(spid)
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.AddRemoteSpidRequest{
		Spid: pbSpid,
	}
	_, err = client.AddRemoteSpid(ctx, request)
	return err
}

func (h *Handler) updateRemoteSpid(pbSpid *pb.Spid) error {
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbSpid.Position))
	if HostIsLocal(ip) {
		spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
		if err != nil {
			return err
		}
		return h.DBManager.UpdateRemoteSpid(spid)
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.UpdateRemoteSpidRequest{
		Spid: pbSpid,
	}
	_, err = client.UpdateRemoteSpid(ctx, request)
	return err
}

func (h *Handler) removeRemoteSpid(pbSpid *pb.Spid) error {
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbSpid.Position))
	if HostIsLocal(ip) {
		spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
		if err != nil {
			return err
		}
		return h.DBManager.RemoveRemoteSpid(spid)
	}
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.RemoveRemoteSpidRequest{
		Spid: pbSpid,
	}
	_, err = client.RemoveRemoteSpid(ctx, request)
	return err
}
