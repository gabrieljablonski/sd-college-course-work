package requestHandling

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
	"time"
)

type localSpidCall func(handler *Handler) (*entities.Spid, error)
type remoteSpidCall func(client pb.SpidHandlerClient, ctx context.Context) (interface{}, error)

func (h *Handler) callSpidGRPC(ip utils.IP, remoteCall remoteSpidCall) (interface{}, error) {
	log.Printf("Making remote call to %s.", ip)
	conn, err := grpc.Dial(ip.String(), grpc.WithInsecure())
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

func (h *Handler) routeSpidCall(ip utils.IP, localCall localSpidCall, remoteCall remoteSpidCall) (*pb.Spid, error) {
	log.Printf("Spid agent: %s", ip)
	if IsHostLocal(ip) {
		log.Printf("Agent is local.")
		spid, err := localCall(h)
		if err != nil || spid == nil {
			return nil, err
		}
		return spid.ToProtoBufferEntity(), nil
	}
	log.Printf("Agent is remote.")
	response, err := h.callSpidGRPC(ip, remoteCall)
	if err != nil {
		return nil, err
	}
	switch t := response.(type) {
	default:
		log.Printf("Invalid type %T", t)
		return nil, fmt.Errorf("invalid response type %T", t)
	case *pb.GetSpidResponse:
		return response.(*pb.GetSpidResponse).Spid, nil
	case *pb.RegisterSpidResponse:
		return response.(*pb.RegisterSpidResponse).Spid, nil
	case *pb.UpdateSpidResponse:
		return response.(*pb.UpdateSpidResponse).Spid, nil
	case *pb.DeleteSpidResponse:
		return response.(*pb.DeleteSpidResponse).Spid, nil
	case *pb.AddRemoteSpidResponse, *pb.UpdateRemoteSpidResponse, *pb.RemoveRemoteSpidResponse:
		return nil, nil
	}
}

func (h *Handler) querySpid(spidID string) (*pb.Spid, error) {
	log.Printf("Querying spid with id %s.", spidID)
	id, err := uuid.Parse(spidID)
	if err != nil {
		return nil, fmt.Errorf("invalid spid id: %s", err)
	}
	ip := h.WhereIsEntity(id)
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return handler.DBManager.QuerySpid(id)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetSpidRequest{
			SpidID: id.String(),
		}
		return client.GetSpidInfo(ctx, request)
	}
	return h.routeSpidCall(ip, localCall, remoteCall)
}

func (h *Handler) registerSpid(batteryLevel uint32, position gps.GlobalPosition) (*pb.Spid, error) {
	log.Printf("Registering spid with battery level %d and position\n%s", batteryLevel, position)
	spid := entities.NewSpid(batteryLevel, position)
	ip := h.WhereIsEntity(spid.ID)
	localCall := func(handler *Handler) (*entities.Spid, error) {
		err := handler.DBManager.RegisterSpid(spid)
		if err != nil {
			return nil, err
		}
		return spid, handler.addRemoteSpid(spid.ToProtoBufferEntity())
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RegisterSpidRequest{
			BatteryLevel: batteryLevel,
			Position:     position.ToProtoBufferEntity(),
		}
		return client.RegisterSpid(ctx, request)
	}
	return h.routeSpidCall(ip, localCall, remoteCall)
}

func (h *Handler) updateSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	log.Printf("Updating spid: %s", spid)
	ip := h.WhereIsEntity(spid.ID)
	localCall := func(handler *Handler) (*entities.Spid, error) {
		err = handler.DBManager.UpdateSpid(spid)
		if err != nil {
			return nil, err
		}
		return nil, handler.updateRemoteSpid(pbSpid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateSpidRequest{
			Spid: pbSpid,
		}
		return client.UpdateSpid(ctx, request)
	}
	_, err = h.routeSpidCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) deleteSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	log.Printf("Updating spid: %s", spid)
	ip := h.WhereIsEntity(spid.ID)
	localCall := func(handler *Handler) (*entities.Spid, error) {
		spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
		if err != nil {
			return nil, err
		}
		err = handler.DBManager.DeleteSpid(spid)
		if err != nil {
			return nil, err
		}
		return nil, handler.removeRemoteSpid(pbSpid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteSpidRequest{
			SpidID: pbSpid.Id,
		}
		return client.DeleteSpid(ctx, request)
	}
	_, err = h.routeSpidCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) addRemoteSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	log.Printf("Adding remote spid: %s", spid)
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbSpid.Position))
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return nil, handler.DBManager.AddRemoteSpid(spid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteSpidRequest{
			Spid: pbSpid,
		}
		return client.AddRemoteSpid(ctx, request)
	}
	_, err = h.routeSpidCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) updateRemoteSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	log.Printf("Updating remote spid: %s", spid)
	ip := h.WhereIsPosition(spid.Position)
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return nil, handler.DBManager.UpdateRemoteSpid(spid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteSpidRequest{
			Spid: pbSpid,
		}
		return client.UpdateRemoteSpid(ctx, request)
	}
	_, err = h.routeSpidCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) removeRemoteSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	log.Printf("Removing remote spid: %s", spid)
	ip := h.WhereIsPosition(spid.Position)
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return nil, handler.DBManager.RemoveRemoteSpid(spid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteSpidRequest{
			Spid: pbSpid,
		}
		return client.RemoveRemoteSpid(ctx, request)
	}
	_, err = h.routeSpidCall(ip, localCall, remoteCall)
	return err
}
