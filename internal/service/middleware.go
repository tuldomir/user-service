package service

import (
	"context"
	"fmt"
	"log"

	"user-service/internal/domain"
	"user-service/pb"
	"user-service/pkg/broker"
	"user-service/pkg/cache"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	userListKey    = "userlist"
	userAddedEvent = "useraddtopic"
)

// CacheMiddleware .
func CacheMiddleware(c cache.Cache) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		if info.FullMethod == "/pb.UserService/Add" ||
			info.FullMethod == "/pb.UserService/Delete" {

			fmt.Println("not list meth, clearing cache")

			if err := c.Clear(ctx, userListKey); err != nil {
				return &pb.ListUsersResponse{},
					status.Errorf(codes.Internal, "cache error %v", err)
			}

			return handler(ctx, req)
		}

		if info.FullMethod != "/pb.UserService/List" {
			return handler(ctx, req)
		}

		list, ok, err := c.Get(ctx, userListKey)
		if err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, "cache error %v", err)
		}

		fmt.Println("got users from cache")

		if ok {
			pbusers := domain.EncodeUserList(list)
			fmt.Println("returnin cache")
			return &pb.ListUsersResponse{Users: pbusers}, nil
		}

		fmt.Println("getting users from real db")

		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}

		r, ok := resp.(*pb.ListUsersResponse)
		if !ok {
			// TODO log error return normal response
			return &pb.ListUsersResponse{},
				status.Error(codes.Internal, "incorrect response type")
		}

		list, err = domain.DecodeUserList(r.Users)
		if err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, err.Error())
		}

		if err = c.Set(ctx, userListKey, list); err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, err.Error())
		}

		return resp, err
	}
}

// KafkaMiddleware .
func KafkaMiddleware(br broker.Broker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		if info.FullMethod != "/pb.UserService/Add" {
			return handler(ctx, req)
		}

		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}

		r, ok := resp.(*pb.AddUserResponse)
		if !ok {
			log.Printf("incorrect response: %v", resp)
			return resp, err
		}

		user, e := domain.DecodeUser(r.User)
		if err != nil {
			log.Printf("cant decode user: %v\n", e)
			return resp, err
		}

		event := &domain.UserEvent{
			EventType: userAddedEvent,
			UID:       user.ID,
			CreatedAt: user.CreatedAt,
		}

		e = br.Publish(ctx, userAddedEvent, event)
		if e != nil {
			log.Printf("cant publish %v: %v\n", event.EventType, e)
			return resp, err
		}

		return resp, err
	}
}
