package service

import (
	"context"
	"user-service/cache"
	"user-service/domain"
	"user-service/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const userListKey = "userlist"

// CacheMiddleware .
func CacheMiddleware(c cache.Cache) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		if info.FullMethod == "/pb.UserService/Add" ||
			info.FullMethod == "/pb.UserService/Delete" {

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

		if ok {
			pbusers := domain.EncodeUserList(list)
			return &pb.ListUsersResponse{Users: pbusers}, nil
		}

		resp, err = handler(ctx, req)
		if err != nil {
			return resp, err
		}

		r, ok := resp.(*pb.ListUsersResponse)
		if !ok {
			return &pb.ListUsersResponse{},
				status.Error(codes.Internal, "incorrect response type")
		}

		list, err = domain.DecodeUserList(r.Users)
		if err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, "decode error %v", err)
		}

		if err = c.Set(ctx, userListKey, list); err != nil {
			return &pb.ListUsersResponse{},
				status.Errorf(codes.Internal, "cache error %v", err)
		}

		return resp, err
	}
}

// KafkaMiddleware .
func KafkaMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		return handler(ctx, req)
	}
}
