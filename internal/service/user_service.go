package service

import (
	"context"
	"time"
	"user-service/internal/mapper"
	"user-service/internal/pb"
	"user-service/internal/repo"
	"user-service/models"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService .
type UserService struct {
	userRepo repo.DB
	pb.UnimplementedUserServiceServer
}

// NewUserService .
func NewUserService(repo repo.DB) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

// Add .
func (s *UserService) Add(
	ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {

	user, err := s.userRepo.Add(ctx, &models.User{
		UID:       uuid.NewString(),
		Email:     req.Email,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return &pb.AddUserResponse{}, status.Errorf(codes.Internal, "cant add user %v", err)
	}

	protoUser := mapper.UserToProto(user)

	return &pb.AddUserResponse{User: protoUser}, nil
}

// Delete .
func (s *UserService) Delete(
	ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {

	err := s.userRepo.Delete(ctx, req.Uuid)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.Internal, "cant delete user %v", err)
	}

	return &empty.Empty{}, err
}

// List .
func (s *UserService) List(
	ctx context.Context, req *empty.Empty) (*pb.ListUsersResponse, error) {

	users, err := s.userRepo.List(ctx)
	if err != nil {
		return &pb.ListUsersResponse{}, status.Errorf(
			codes.Internal, "cant get users %v", err)
	}

	protoUsers := mapper.UserToProtoList(users)
	return &pb.ListUsersResponse{Users: protoUsers}, nil
}
