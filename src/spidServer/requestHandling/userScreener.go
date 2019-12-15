package requestHandling

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"spidServer/db"
	"spidServer/entities"
	"spidServer/gps"
	pb "spidServer/requestHandling/protoBuffers"
	"spidServer/utils"
)

type localUserCall func(handler *Handler) (*entities.User, error)
type remoteUserCall func(client pb.SpidHandlerClient, ctx context.Context) (interface{}, error)

func (h *Handler) callUserGRPC(ip utils.IP, remoteCall remoteUserCall) (interface{}, error) {
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

func (h *Handler) routeUserCall(targetServerNumber int, localCall localUserCall, remoteCall remoteUserCall) (pbUser *pb.User, err error) {
	log.Printf("Target server number: %d", targetServerNumber)
	if targetServerNumber == h.ServerNumber {
		log.Printf("Agent is local (%s).", h.IPMap[targetServerNumber])
		user, err := localCall(h)
		if err != nil || user == nil {
			return nil, err
		}
		return user.ToProtoBufferEntity()
	}
	ips := h.getClosestHost(targetServerNumber)
	log.Printf("Agent is remote: %s.", ips)
	var response interface{}
	for _, ip := range ips {
		response, err = h.callUserGRPC(ip, remoteCall)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to complete request to to any servers in %#v", ips)
	}
	switch t := response.(type) {
	default:
		log.Printf("Invalid type %T", t)
		return nil, fmt.Errorf("invalid response type %T", t)
	case *pb.GetUserResponse:
		return response.(*pb.GetUserResponse).User, nil
	case *pb.RegisterUserResponse:
		return response.(*pb.RegisterUserResponse).User, nil
	case *pb.UpdateUserResponse:
		return response.(*pb.UpdateUserResponse).User, nil
	case *pb.DeleteUserResponse:
		return response.(*pb.DeleteUserResponse).User, nil
	case *pb.AddRemoteUserResponse, *pb.UpdateRemoteUserResponse, *pb.RemoveRemoteUserResponse:
		return nil, nil
	}
}

func (h *Handler) queryUser(userID string) (*pb.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", err)
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		return handler.ConsensusManager.DBManager.QueryUser(id)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetUserRequest{
			UserID: id.String(),
		}
		log.Printf("Sending GetUserInfo request: %s.", request)
		return client.GetUserInfo(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(id)
	return h.routeUserCall(targetServerNumber, localCall, remoteCall)
}

func (h *Handler) registerUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		action := db.WriteAction{
			Location:   db.Local,
			EntityType: db.User,
			Type:       db.Register,
			UserEntity: user,
		}
		err := handler.ConsensusManager.PutCommand(action)
		if err != nil {
			return nil, err
		}
		pbUser, _ := user.ToProtoBufferEntity()
		return user, handler.addRemoteUser(pbUser)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RegisterUserRequest{
			User: pbUser,
		}
		log.Printf("Sending RegisterUser request: %s.", request)
		return client.RegisterUser(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(user.ID)
	_, err = h.routeUserCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) updateUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		action := db.WriteAction{
			Location:   db.Local,
			EntityType: db.User,
			Type:       db.Update,
			UserEntity: user,
		}
		err := handler.ConsensusManager.PutCommand(action)
		if err != nil {
			return nil, err
		}
		return nil, handler.updateRemoteUser(pbUser)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateUserRequest{
			User: pbUser,
		}
		log.Printf("Sending UpdateUser request: %s.", request)
		return client.UpdateUser(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(user.ID)
	_, err = h.routeUserCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) deleteUser(userID string) (*pb.User, error) {
	pbUser, err := h.queryUser(userID)
	if err != nil {
		return nil, err
	}
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return nil, err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		action := db.WriteAction{
			Location:   db.Local,
			EntityType: db.User,
			Type:       db.Delete,
			UserEntity: user,
		}
		err := handler.ConsensusManager.PutCommand(action)
		if err != nil {
			return nil, err
		}
		return nil, handler.removeRemoteUser(userID)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteUserRequest{
			UserID: userID,
		}
		log.Printf("Sending DeleteUser request: %s.", request)
		return client.DeleteUser(ctx, request)
	}
	targetServerNumber := h.WhereIsEntity(user.ID)
	_, err = h.routeUserCall(targetServerNumber, localCall, remoteCall)
	return pbUser, err
}

func (h *Handler) addRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		action := db.WriteAction{
			Location:   db.Remote,
			EntityType: db.User,
			Type:       db.Add,
			UserEntity: user,
		}
		return nil, handler.ConsensusManager.PutCommand(action)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteUserRequest{
			User: pbUser,
		}
		log.Printf("Sending AddRemoteUser request: %s.", request)
		return client.AddRemoteUser(ctx, request)
	}
	targetServerNumber := h.WhereIsPosition(user.Position)
	_, err = h.routeUserCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) updateRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		action := db.WriteAction{
			Location:   db.Remote,
			EntityType: db.User,
			Type:       db.Update,
			UserEntity: user,
		}
		return nil, handler.ConsensusManager.PutCommand(action)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteUserRequest{
			User: pbUser,
		}
		log.Printf("Sending UpdateRemoteUser request: %s.", request)
		return client.UpdateRemoteUser(ctx, request)
	}
	targetServerNumber := h.WhereIsPosition(user.Position)
	_, err = h.routeUserCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) removeRemoteUser(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	user, err := h.ConsensusManager.DBManager.QueryRemoteUser(uid)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		action := db.WriteAction{
			Location:   db.Remote,
			EntityType: db.User,
			Type:       db.Remove,
			UserEntity: user,
		}
		return nil, handler.ConsensusManager.PutCommand(action)
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteUserRequest{
			UserID: userID,
		}
		log.Printf("Sending RemoveRemoteUser request: %s.", request)
		return client.RemoveRemoteUser(ctx, request)
	}
	targetServerNumber := h.WhereIsPosition(user.Position)
	_, err = h.routeUserCall(targetServerNumber, localCall, remoteCall)
	return err
}

func (h *Handler) getRemoteSpids(pbPosition *pb.GlobalPosition) (string, error) {
	position, err := gps.FromProtoBufferEntity(pbPosition)
	if err != nil {
		return "", err
	}
	targetServerNumber := h.WhereIsPosition(position)
	if targetServerNumber == h.ServerNumber {
		log.Printf("Agent is local (%s).", h.IPMap[targetServerNumber])
		marshaledSpids, err := json.Marshal(h.ConsensusManager.DBManager.GetRemoteSpids())
		return string(marshaledSpids), err
	}
	remoteCall := func (client pb.SpidHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetRemoteSpidsRequest{
			Position: pbPosition,
		}
		log.Printf("Sending RemoveRemoteSpid request: %s.", request)
		return client.GetRemoteSpids(ctx, request)
	}
	ips := h.getClosestHost(targetServerNumber)
	log.Printf("Agent is remote: %s.", ips)
	var response interface{}
	for _, ip := range ips {
		response, err = h.callUserGRPC(ip, remoteCall)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to complete request to to any servers in %#v", ips)
	}
	return response.(*pb.GetRemoteSpidsResponse).MarshaledSpids, nil
}
