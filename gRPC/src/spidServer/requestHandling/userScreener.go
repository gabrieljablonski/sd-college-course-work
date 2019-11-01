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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
		return user.ToProtoBufferEntity(), nil
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
	log.Printf("Querying user with id %s.", userID)
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", err)
	}
	ip := h.WhereIsEntity(id)
	if IsHostLocal(ip) {
		user, err := h.DBManager.QueryUser(id)
		if err != nil {
			return nil, err
		}
		return user.ToProtoBufferEntity(), err
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.GetUserRequest{
			UserID: id.String(),
		}
		return client.GetUserInfo(ctx, request)
	}
	response, err := h.callUserGRPC(ip, remoteCall)
	if err != nil {
		return nil, err
	}
	return response.(*pb.GetUserResponse).User, nil
}

func (h *Handler) registerUser(name string, position gps.GlobalPosition) (*pb.User, error) {
	log.Printf("Registering user with name `%s` and position\n%s", name, position)
	user := entities.NewUser(name, position)
	ip := h.WhereIsEntity(user.ID)
	localCall := func(handler *Handler) (*entities.User, error) {
		err := handler.DBManager.RegisterUser(user)
		if err != nil {
			return nil, err
		}
		return user, handler.addRemoteUser(user.ToProtoBufferEntity())
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RegisterUserRequest{
			Name: name,
		}
		return client.RegisterUser(ctx, request)
	}
	return h.routeUserCall(ip, localCall, remoteCall)
}

func (h *Handler) updateUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	log.Printf("Updating user: %s", user)
	ip := h.WhereIsEntity(user.ID)
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
		return client.UpdateUser(ctx, request)
	}
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) deleteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	log.Printf("Deleting user: %s", user)
	ip := h.WhereIsEntity(user.ID)
	localCall := func(handler *Handler) (*entities.User, error) {
		err = handler.DBManager.DeleteUser(user)
		if err != nil {
			return nil, err
		}
		return nil, handler.removeRemoteUser(pbUser)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteUserRequest{
			UserID: pbUser.Id,
		}
		return client.DeleteUser(ctx, request)
	}
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) addRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	log.Printf("Adding remote user: %s", user)
	ip := h.WhereIsPosition(user.Position)
	localCall := func(handler *Handler) (*entities.User, error) {
		return nil, handler.DBManager.AddRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteUserRequest{
			User: pbUser,
		}
		return client.AddRemoteUser(ctx, request)
	}
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) updateRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	log.Printf("Updating remote user: %s", user)
	ip := h.WhereIsPosition(user.Position)
	localCall := func(handler *Handler) (*entities.User, error) {
		return nil, handler.DBManager.UpdateRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteUserRequest{
			User: pbUser,
		}
		return client.UpdateRemoteUser(ctx, request)
	}
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}

func (h *Handler) removeRemoteUser(pbUser *pb.User) error {
	user, err := entities.UserFromProtoBufferEntity(pbUser)
	if err != nil {
		return err
	}
	log.Printf("Deleting remote user: %s", user)
	ip := h.WhereIsPosition(user.Position)
	localCall := func(handler *Handler) (*entities.User, error) {
		return nil, handler.DBManager.RemoveRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteUserRequest{
			User: pbUser,
		}
		return client.RemoveRemoteUser(ctx, request)
	}
	_, err = h.routeUserCall(ip, localCall, remoteCall)
	return err
}