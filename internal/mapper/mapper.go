package mapper

import (
	"user-service/internal/pb"
	"user-service/models"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserToProto .
func UserToProto(user *models.User) *pb.User {
	createdAt := timestamppb.New(user.CreatedAt)

	return &pb.User{
		Uuid:      user.ID.String(),
		Email:     user.Email,
		CreatedAt: createdAt,
	}
}

// ProtoToUser .
func ProtoToUser(user *pb.User) (*models.User, error) {
	id, err := uuid.Parse(user.Uuid)
	if err != nil {
		return nil, err
	}

	createdAt := user.CreatedAt.AsTime()

	return &models.User{
		ID:        id,
		Email:     user.Email,
		CreatedAt: createdAt,
	}, nil
}

// UserToProtoList .
func UserToProtoList(users []*models.User) []*pb.User {
	protoUsers := make([]*pb.User, len(users))

	for i, u := range users {
		pbUser := UserToProto(u)
		protoUsers[i] = pbUser
	}

	return protoUsers
}

// ProtoToUserList .
func ProtoToUserList(pbusers []*pb.User) ([]*models.User, error) {
	users := make([]*models.User, len(pbusers))

	for i, u := range pbusers {
		user, err := ProtoToUser(u)
		if err != nil {
			return nil, err
		}

		users[i] = user
	}

	return users, nil
}
