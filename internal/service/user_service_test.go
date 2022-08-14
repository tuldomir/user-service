package service

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"
	"user-service/internal/mapper"
	"user-service/internal/pb"
	mockdb "user-service/internal/repo/mock"
	"user-service/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

const bufSize = 1024 * 1024

func bufListener(
	serv pb.UserServiceServer,
	quit chan struct{}, opts ...grpc.ServerOption) func(context.Context, string) (net.Conn, error) {

	lis := bufconn.Listen(bufSize)

	var s *grpc.Server
	s = grpc.NewServer(opts...)

	pb.RegisterUserServiceServer(s, serv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	go func() {
		<-quit
		s.Stop()
	}()

	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func Test_UserService_Add(t *testing.T) {

	user := &models.User{
		UID:       uuid.NewString(),
		Email:     "foo@mail.test",
		CreatedAt: time.Now().UTC(),
	}

	pbUser := mapper.UserToProto(user)
	assert.NotEmpty(t, pbUser)

	testCases := []struct {
		name        string
		mockStubs   func(*mockdb.MockDB)
		checkResult func(*pb.AddUserResponse, error)
	}{
		{
			name: "ok",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					Add(gomock.Any(), gomock.Any()).
					Times(1).Return(user, nil)
			},
			checkResult: func(res *pb.AddUserResponse, err error) {
				assert.NoErrorf(t, err, "Add method failed %v", err)
				assert.NotNil(t, res)
				assert.Equalf(t, pbUser, res.User,
					"expected : %v, actual: %v", pbUser, res.User)

			},
		},
		{
			name: "dbError",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					Add(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, errors.New("db error"))
			},
			checkResult: func(res *pb.AddUserResponse, err error) {
				expectedErr := status.Errorf(
					codes.Internal, "cant add user %v", errors.New("db error"))

				assert.ErrorIsf(
					t, expectedErr, err, "expected: %v, actual: %v", expectedErr, err)
				assert.Empty(t, res)

			},
		},
	}

	quit := make(chan struct{})
	defer close(quit)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mockdb.NewMockDB(ctrl)

	serv := NewUserService(db)

	bufDialer := bufListener(serv, quit)

	conn, err := grpc.DialContext(
		context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure())

	assert.NoErrorf(t, err, "failed to dial bufnet %v", err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {

			v.mockStubs(db)

			res, err := client.Add(
				context.Background(), &pb.AddUserRequest{Email: user.Email})

			v.checkResult(res, err)
		})
	}
}

func Test_UserService_Delete(t *testing.T) {
	uid := uuid.NewString()

	testCases := []struct {
		name      string
		mockStubs func(*mockdb.MockDB)
		id        string
		err       error
	}{
		{
			name: "ok",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					Delete(gomock.Any(), gomock.Eq(uid)).
					Times(1).Return(nil)
			},
			id:  uid,
			err: nil,
		},
		{
			name: "dbError",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					Delete(gomock.Any(), gomock.Eq(uid)).
					Times(1).Return(errors.New("db error"))
			},
			id: uid,
			err: status.Errorf(
				codes.Internal, "cant delete user %v", errors.New("db error")),
		},
	}

	quit := make(chan struct{})
	defer close(quit)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mockdb.NewMockDB(ctrl)

	serv := NewUserService(db)

	bufDialer := bufListener(serv, quit)

	conn, err := grpc.DialContext(
		context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure())

	assert.NoErrorf(t, err, "failed to dial bufnet %v", err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {

			v.mockStubs(db)

			_, err := client.Delete(
				context.Background(), &pb.DeleteUserRequest{Uuid: uid})

			assert.ErrorIsf(t, err, v.err, "expected %v, actual: %v", v.err, err)
		})
	}
}

func Test_UserService_List(t *testing.T) {

	user := &models.User{
		UID:       uuid.NewString(),
		Email:     "foo@mail.test",
		CreatedAt: time.Now().UTC(),
	}

	users := []*models.User{user}

	testCases := []struct {
		name        string
		mockStubs   func(*mockdb.MockDB)
		checkResult func(*pb.ListUsersResponse, error)
	}{
		{
			name: "ok",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					List(gomock.Any()).
					Times(1).Return(users, nil)
			},
			checkResult: func(res *pb.ListUsersResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				list, err := mapper.ProtoToUserList(res.Users)
				assert.NoError(t, err)
				assert.NotEmpty(t, list)

				assert.Equal(t, users, list)
			},
		},
		{
			name: "dbError",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					List(gomock.Any()).
					Times(1).Return(nil, errors.New("db error"))
			},
			checkResult: func(res *pb.ListUsersResponse, err error) {
				assert.NotNil(t, err)

				expectedErr := status.Errorf(
					codes.Internal, "cant get users %v", errors.New("db error"))

				assert.ErrorIsf(
					t, expectedErr, err, "expected: %v, actual: %v", expectedErr, err)

				assert.Nil(t, res)
			},
		},
	}

	quit := make(chan struct{})
	defer close(quit)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mockdb.NewMockDB(ctrl)

	serv := NewUserService(db)

	bufDialer := bufListener(serv, quit)

	conn, err := grpc.DialContext(
		context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure())

	assert.NoErrorf(t, err, "failed to dial bufnet %v", err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {

			v.mockStubs(db)

			res, err := client.List(
				context.Background(), &emptypb.Empty{})

			v.checkResult(res, err)
		})
	}
}
