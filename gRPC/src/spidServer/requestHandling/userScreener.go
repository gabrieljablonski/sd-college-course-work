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

type remoteUserCall func(client pb.UserHandlerClient, ctx context.Context) (interface{}, error)

func (h *Handler) callUserGRPC(ip utils.IP, remoteCall remoteUserCall) (interface{}, error) {
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
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

func (h *Handler) queryUser(userID string) (*pb.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", err)
	}
	ip := h.WhereIsEntity(id)
	if HostIsLocal(ip) {
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
	return response.(pb.GetUserResponse).User, nil
}

func (h *Handler) registerUser(name string, position *pb.GlobalPosition) (*pb.User, error) {
	user := entities.NewUser(name, gps.FromProtoBufferEntity(position))
	ip := h.WhereIsEntity(user.ID)
	if HostIsLocal(ip) {
		err := h.DBManager.RegisterUser(user)
		if err != nil {
			return nil, err
		}
		err = h.addRemoteUser(user.ToProtoBufferEntity())
		return user.ToProtoBufferEntity(), err
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RegisterUserRequest{
			Name: name,
		}
		return client.RegisterUser(ctx, request)
	}
	response, err := h.callUserGRPC(ip, remoteCall)
	if err != nil {
		return nil, err
	}
	return response.(pb.RegisterUserResponse).User, nil
}

func (h *Handler) updateUser(pbUser *pb.User) error {
	id, err := uuid.Parse(pbUser.Id)
	if err != nil {
		return err
	}
	ip := h.WhereIsEntity(id)
	if HostIsLocal(ip) {
		user, err := entities.UserFromProtoBufferEntity(pbUser)
		if err != nil {
			return err
		}
		err = h.DBManager.UpdateUser(user)
		if err != nil {
			return err
		}
		return h.updateRemoteUser(pbUser)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateUserRequest{
			User: pbUser,
		}
		return client.UpdateUser(ctx, request)
	}
	_, err = h.callUserGRPC(ip, remoteCall)
	return err
}

func (h *Handler) deleteUser(pbUser *pb.User) error {
	id, err := uuid.Parse(pbUser.Id)
	if err != nil {
		return err
	}
	ip := h.WhereIsEntity(id)
	if HostIsLocal(ip) {
		user, err := entities.UserFromProtoBufferEntity(pbUser)
		if err != nil {
			return err
		}
		err = h.DBManager.DeleteUser(user)
		if err != nil {
			return err
		}
		return h.removeRemoteUser(pbUser)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.DeleteUserRequest{
			UserID: pbUser.Id,
		}
		return client.DeleteUser(ctx, request)
	}
	_, err = h.callUserGRPC(ip, remoteCall)
	return err
}

func (h *Handler) addRemoteUser(pbUser *pb.User) error {
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbUser.Position))
	if HostIsLocal(ip) {
		user, err := entities.UserFromProtoBufferEntity(pbUser)
		if err != nil {
			return err
		}
		return h.DBManager.AddRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.AddRemoteUserRequest{
			User: pbUser,
		}
		return client.AddRemoteUser(ctx, request)
	}
	_, err := h.callUserGRPC(ip, remoteCall)
	return err
}

func (h *Handler) updateRemoteUser(pbUser *pb.User) error {
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbUser.Position))
	if HostIsLocal(ip) {
		user, err := entities.UserFromProtoBufferEntity(pbUser)
		if err != nil {
			return err
		}
		return h.DBManager.UpdateRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.UpdateRemoteUserRequest{
			User: pbUser,
		}
		return client.UpdateRemoteUser(ctx, request)
	}
	_, err := h.callUserGRPC(ip, remoteCall)
	return err
}

func (h *Handler) removeRemoteUser(pbUser *pb.User) error {
	ip := h.WhereIsPosition(gps.FromProtoBufferEntity(pbUser.Position))
	if HostIsLocal(ip) {
		user, err := entities.UserFromProtoBufferEntity(pbUser)
		if err != nil {
			return err
		}
		return h.DBManager.RemoveRemoteUser(user)
	}
	remoteCall := func (client pb.UserHandlerClient, ctx context.Context) (interface{}, error) {
		request := &pb.RemoveRemoteUserRequest{
			User: pbUser,
		}
		return client.RemoveRemoteUser(ctx, request)
	}
	_, err := h.callUserGRPC(ip, remoteCall)
	return err
}