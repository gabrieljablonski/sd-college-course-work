package requestHandling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
	"time"
)

type remoteSpidCall func(client pb.SpidHandlerClient, ctx context.Context) (interface{}, error)

func (h *Handler) callSpidGRPC(ip utils.IP, remoteCall remoteSpidCall) (interface{}, error) {
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	client := pb.NewSpidHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return remoteCall(client, ctx)
}

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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetSpidRequest{
			SpidID: id.String(),
		}
		return client.GetSpidInfo(ctx, request)
	}
	response, err := h.callSpidGRPC(ip, remoteCall)
	if err != nil {
		return nil, err
	}
	return response.(pb.GetSpidResponse).Spid, nil
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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RegisterSpidRequest{
			BatteryLevel: batteryLevel,
			Position:     position.ToProtoBufferEntity(),
		}
		return client.RegisterSpid(ctx, request)
	}
	response, err := h.callSpidGRPC(ip, remoteCall)
	if err != nil {
		return nil, err
	}
	return response.(pb.RegisterSpidResponse).Spid, nil
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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateSpidRequest{
			Spid: pbSpid,
		}
		return client.UpdateSpid(ctx, request)
	}
	_, err = h.callSpidGRPC(ip, remoteCall)
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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteSpidRequest{
			SpidID: pbSpid.Id,
		}
		return client.DeleteSpid(ctx, request)
	}
	_, err = h.callSpidGRPC(ip, remoteCall)
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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteSpidRequest{
			Spid: pbSpid,
		}
		return client.AddRemoteSpid(ctx, request)
	}
	_, err := h.callSpidGRPC(ip, remoteCall)
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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteSpidRequest{
			Spid: pbSpid,
		}
		return client.UpdateRemoteSpid(ctx, request)
	}
	_, err := h.callSpidGRPC(ip, remoteCall)
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
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteSpidRequest{
			Spid: pbSpid,
		}
		return client.RemoveRemoteSpid(ctx, request)
	}
	_, err := h.callSpidGRPC(ip, remoteCall)
	return err
}
