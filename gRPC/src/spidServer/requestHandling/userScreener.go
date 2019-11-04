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

type localUserCall func(handler *Handler) (*entities.User, error)
type remoteUserCall func(client pb.UserHandlerClient, ctx context.Context) (interface{}, error)

func (h *Handler) callUserGRPC(ip utils.IP, remoteCall remoteUserCall) (interface{}, error) {
	log.Printf("Making remote call to %s.", ip)
	conn, err := grpc.Dial(ip.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()
	return remoteCall(client, ctx)
}

func (h *Handler) routeUserCall(ip utils.IP, localCall localUserCall, remoteCall remoteUserCall) (*pb.User, error) {
	log.Printf("User agent: %s", ip)
	if IsHostLocal(ip) {
		log.Printf("Agent is local.")
		user, err := localCall(h)
		if err != nil || user == nil {
			return nil, err
		}
		return user.ToProtoBufferEntity()
	}
	log.Printf("Agent is remote.")
	response, err := h.callUserGRPC(ip, remoteCall)
	if err != nil {
		return nil, err
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
		return handler.DBManager.QueryUser(id)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetUserRequest{
			UserID: id.String(),
		}
		log.Printf("Sending GetUserInfo request: %s.", request)
		return client.GetUserInfo(ctx, request)
	}
	ip := h.WhereIsEntity(id)
	return h.routeUserCall(ip, localCall, remoteCall)
}

func (h *Handler) registerUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		err := handler.DBManager.RegisterUser(user)
		if err != nil {
			return nil, err
		}
		pbUser, _ := user.ToProtoBufferEntity()
		return user, handler.addRemoteUser(pbUser)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RegisterUserRequest{
			User: pbUser,
		}
		log.Printf("Sending RegisterUser request: %s.", request)
		return client.RegisterUser(ctx, request)
	}
	ip := h.WhereIsEntity(user.ID)
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) updateUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		err = handler.DBManager.UpdateUser(user)
		if err != nil {
			return nil, err
		}
		return nil, handler.updateRemoteUser(pbUser)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateUserRequest{
			User: pbUser,
		}
		log.Printf("Sending UpdateUser request: %s.", request)
		return client.UpdateUser(ctx, request)
	}
	ip := h.WhereIsEntity(user.ID)
	_, err = h.routeUserCall(ip, localCall, remoteCall)
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
		err = handler.DBManager.DeleteUser(user.ID)
		if err != nil {
			return nil, err
		}
		return nil, handler.removeRemoteUser(userID)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteUserRequest{
			UserID: userID,
		}
		log.Printf("Sending DeleteUser request: %s.", request)
		return client.DeleteUser(ctx, request)
	}
	ip := h.WhereIsEntity(user.ID)
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return pbUser, err
}

func (h *Handler) addRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		return nil, handler.DBManager.AddRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteUserRequest{
			User: pbUser,
		}
		log.Printf("Sending AddRemoteUser request: %s.", request)
		return client.AddRemoteUser(ctx, request)
	}
	ip := h.WhereIsPosition(user.Position)
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) updateRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		return nil, handler.DBManager.UpdateRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteUserRequest{
			User: pbUser,
		}
		log.Printf("Sending UpdateRemoteUser request: %s.", request)
		return client.UpdateRemoteUser(ctx, request)
	}
	ip := h.WhereIsPosition(user.Position)
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) removeRemoteUser(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	localCall := func(handler *Handler) (*entities.User, error) {
		return nil, handler.DBManager.RemoveRemoteUser(uid)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteUserRequest{
			UserID: userID,
		}
		log.Printf("Sending RemoveRemoteUser request: %s.", request)
		return client.RemoveRemoteUser(ctx, request)
	}
	user, err := h.queryUser(userID)
	if err != nil {
		return err
	}
	pbPosition, _ := gps.FromProtoBufferEntity(user.Position)
	ip := h.WhereIsPosition(pbPosition)
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}