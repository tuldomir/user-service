package service

import (
	"context"
	"log"

	"user-service/internal/mapper"
	"user-service/internal/pb"
	"user-service/models"
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

		if info.FullMethod != "/pb.UserService/List" {
			log.Println("not list meth, clearing cache")

			if err := c.Clear(ctx, userListKey); err != nil {
				log.Printf("error clearing cache: %v\n", err)
			}

			return handler(ctx, req)
		}

		users, ok, err := c.Get(ctx, userListKey)
		if err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, "cache error %v", err)
		}

		// FROM CACHE .
		if ok {
			log.Println("got users from cache")
			pbusers := mapper.UserToProtoList(users)
			return &pb.ListUsersResponse{Users: pbusers}, nil
		}

		// FROM DB .
		log.Println("getting users from real db")

		resp, err = handler(ctx, req)
		if err != nil {
			// err nothing to cache
			return resp, err
		}

		r, ok := resp.(*pb.ListUsersResponse)
		if !ok {
			log.Printf("cant cast response: %v\n", r)
			return &pb.ListUsersResponse{},
				status.Error(codes.Internal, "incorrect response type")
		}

		users, err = mapper.ProtoToUserList(r.Users)
		if err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, err.Error())
		}

		if err = c.Set(ctx, userListKey, users); err != nil {
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

		user, e := mapper.ProtoToUser(r.User)
		if err != nil {
			log.Printf("cant decode user: %v\n", e)
			return resp, err
		}

		event := &models.UserEvent{
			EventType: userAddedEvent,
			UID:       user.UID,
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
