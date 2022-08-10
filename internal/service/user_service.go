package service

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"
	"user-service/pb"
	"user-service/pkg/cache"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService .
type UserService struct {
	userRepo  repo.DB
	userCache cache.Cache
	pb.UnimplementedUserServiceServer
}

// NewUserService .
func NewUserService(repo repo.DB, cache cache.Cache) *UserService {
	return &UserService{
		userRepo:  repo,
		userCache: cache,
	}
}

// Add .
func (s *UserService) Add(
	ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return &pb.AddUserResponse{}, err
	}

	fmt.Println("line 40", req.Email, id.String())

	user, err := s.userRepo.Add(ctx, &domain.User{
		ID:        id,
		Email:     req.Email,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return &pb.AddUserResponse{}, status.Errorf(codes.Internal, "cant add user %v", err)
	}

	protoUser := domain.EncodeUser(user)

	// TODO publish to kafka useradded event in middleware/interceptor

	return &pb.AddUserResponse{User: protoUser}, nil
}

// Delete .
func (s *UserService) Delete(
	ctx context.Context, req *pb.DeleteUserRequest) (*empty.Empty, error) {

	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.Internal, "cant parse id %v", err)
	}

	err = s.userRepo.Delete(ctx, id)
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

	protoUsers := domain.EncodeUserList(users)
	return &pb.ListUsersResponse{Users: protoUsers}, nil
}
