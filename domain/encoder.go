package domain

import (
	"user-service/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EncodeUser .
func EncodeUser(user *User) *pb.User {
	createdAt := timestamppb.New(user.CreatedAt)

	return &pb.User{
		Uuid:      user.ID.String(),
		Email:     user.Email,
		CreatedAt: createdAt,
	}
}

// DecodeUser .
func DecodeUser(user *pb.User) (*User, error) {
	id, err := uuid.Parse(user.Uuid)
	if err != nil {
		return nil, err
	}

	createdAt := user.CreatedAt.AsTime()

	return &User{
		ID:        id,
		Email:     user.Email,
		CreatedAt: createdAt,
	}, nil
}

// EncodeUserList .
func EncodeUserList(users []*User) []*pb.User {
	protoUsers := make([]*pb.User, len(users))

	for i, u := range users {
		pbUser := EncodeUser(u)
		protoUsers[i] = pbUser
	}

	return protoUsers
}

// DecodeUserList .
func DecodeUserList(pbusers []*pb.User) ([]*User, error) {
	users := make([]*User, len(pbusers))

	for i, u := range pbusers {
		user, err := DecodeUser(u)
		if err != nil {
			return nil, err
		}

		users[i] = user
	}

	return users, nil
}
