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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.GetUserRequest{
		UserID: id.String(),
	}
	response, err := client.GetUserInfo(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.User, err
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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.RegisterUserRequest{
		Name: name,
	}
	response, err := client.RegisterUser(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.User, nil
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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.UpdateUserRequest{
		User: pbUser,
	}
	_, err = client.UpdateUser(ctx, request)
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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.DeleteUserRequest{
		UserID: pbUser.Id,
	}
	_, err = client.DeleteUser(ctx, request)
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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.AddRemoteUserRequest{
		User: pbUser,
	}
	_, err = client.AddRemoteUser(ctx, request)
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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.UpdateRemoteUserRequest{
		User: pbUser,
	}
	_, err = client.UpdateRemoteUser(ctx, request)
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
	conn, err := grpc.Dial(ip.ToString(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewUserHandlerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := &pb.RemoveRemoteUserRequest{
		User: pbUser,
	}
	_, err = client.RemoveRemoteUser(ctx, request)
	return err
}