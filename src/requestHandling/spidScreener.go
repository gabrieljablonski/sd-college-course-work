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
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()
	return remoteCall(client, ctx)
}

func (h *Handler) routeSpidCall(targetServerNumber int, localCall localSpidCall, remoteCall remoteSpidCall) (*pb.Spid, error) {
	log.Printf("Target server number: %d", targetServerNumber)
	if targetServerNumber == h.ServerNumber {
		log.Printf("Agent is local (%s).", h.IPMap[targetServerNumber])
		spid, err := localCall(h)
		if err != nil || spid == nil {
			return nil, err
		}
		return spid.ToProtoBufferEntity()
	}
	ip := h.getClosestHost(targetServerNumber)
	log.Printf("Agent is remote: %s.", ip)
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
	id, err := uuid.Parse(spidID)
	if err != nil {
		return nil, fmt.Errorf("invalid spid id: %s", err)
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return handler.DBManager.QuerySpid(id)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetSpidRequest{
			SpidID: id.String(),
		}
		log.Printf("Sending GetSpidInfo request: %s.", request)
		return client.GetSpidInfo(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(id)
	return h.routeSpidCall(targetServerNumber, localCall, remoteCall)
}

func (h *Handler) registerSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		err := handler.DBManager.RegisterSpid(spid)
		if err != nil {
			return nil, err
		}
		pbSpid, err := spid.ToProtoBufferEntity()
		if err != nil {
			return nil, err
		}
		return spid, handler.addRemoteSpid(pbSpid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		// already checked on NewSpid() call, no need to check again
		request := &pb.RegisterSpidRequest{
			Spid: pbSpid,
		}
		log.Printf("Sending RegisterSpid request: %s.", request)
		return client.RegisterSpid(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(spid.ID)
	_, err = h.routeSpidCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) updateSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		oldSpid, err := handler.DBManager.QuerySpid(spid.ID)
		if err != nil {
			return nil, err
		}
		oldPosition := oldSpid.Position
		newPosition := spid.Position
		if handler.WhereIsPosition(oldPosition) != handler.WhereIsPosition(newPosition) {
			// crossed over a boundary
			err = handler.removeRemoteSpid(oldSpid.ID.String())
			if err != nil {
				return nil, err
			}
			err = handler.addRemoteSpid(pbSpid)
			if err != nil {
				return nil, err
			}
		} else {
			err = handler.updateRemoteSpid(pbSpid)
			if err != nil {
				return nil, err
			}
		}
		return nil, handler.DBManager.UpdateSpid(spid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateSpidRequest{
			Spid: pbSpid,
		}
		log.Printf("Sending UpdateSpid request: %s.", request)
		return client.UpdateSpid(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(spid.ID)
	_, err = h.routeSpidCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) deleteSpid(spidID string) (*pb.Spid, error) {
	pbSpid, err := h.querySpid(spidID)
	if err != nil {
		return nil, err
	}
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return nil, err
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		err = handler.DBManager.DeleteSpid(spid.ID)
		if err != nil {
			return nil, err
		}
		return nil, handler.removeRemoteSpid(spidID)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteSpidRequest{
			SpidID: spidID,
		}
		log.Printf("Sending DeleteSpid request: %s.", request)
		return client.DeleteSpid(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(spid.ID)
	_, err = h.routeSpidCall(targetServerNumber, localCall, remoteCall)
	return pbSpid, err
}

func (h *Handler) addRemoteSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return nil, handler.DBManager.AddRemoteSpid(spid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteSpidRequest{
			Spid: pbSpid,
		}
		log.Printf("Sending AddRemoteSpid request: %s.", request)
		return client.AddRemoteSpid(ctx, request)
	}
	targetServerNumber := h.WhereIsPosition(spid.Position)
	_, err = h.routeSpidCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) updateRemoteSpid(pbSpid *pb.Spid) error {
	spid, err := entities.SpidFromProtoBufferEntity(pbSpid)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return nil, handler.DBManager.UpdateRemoteSpid(spid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteSpidRequest{
			Spid: pbSpid,
		}
		log.Printf("Sending UpdateRemoteSpid request: %s.", request)
		return client.UpdateRemoteSpid(ctx, request)
	}
	targetServerNumber := h.WhereIsPosition(spid.Position)
	_, err = h.routeSpidCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) removeRemoteSpid(spidID string) error {
	uid, err := uuid.Parse(spidID)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.Spid, error) {
		return nil, handler.DBManager.RemoveRemoteSpid(uid)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteSpidRequest{
			SpidID: spidID,
		}
		log.Printf("Sending RemoveRemoteSpid request: %s.", request)
		return client.RemoveRemoteSpid(ctx, request)
	}
	spid, err := h.querySpid(spidID)
	if err != nil {
		return err
	}
	pbPosition, _ := gps.FromProtoBufferEntity(spid.Position)
	targetServerNumber := h.WhereIsPosition(pbPosition)
	_, err = h.routeSpidCall(targetServerNumber, localCall, remoteCall)
	return err
}

