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
)

const bufSize = 1024 * 1024

func bufListener(
	serv pb.UserServiceServer,
	quit chan struct{}) func(context.Context, string) (net.Conn, error) {

	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
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
	id, err := uuid.NewUUID()
	assert.NoError(t, err)

	user := &models.User{
		ID:        id,
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
					Add(gomock.Any(),
						gomock.Eq(user)).
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
					Add(gomock.Any(), gomock.Eq(user)).
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

	// fakecache := cache.NewFakeCache()

	serv := NewUserService(db, nil)

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
	id, err := uuid.NewUUID()
	assert.NoError(t, err)

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
					Delete(gomock.Any(), gomock.Eq(id)).
					Times(1).Return(nil)
			},
			id:  id.String(),
			err: nil,
		},
		{
			name: "uuidError",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Times(0)
			},
			id: "gibberish",
			err: status.Errorf(
				codes.Internal, "cant parse id %v", errors.New("invalid UUID length: 9")),
		},
		{
			name: "dbError",
			mockStubs: func(db *mockdb.MockDB) {
				db.EXPECT().
					Delete(gomock.Any(), gomock.Eq(id)).
					Times(1).Return(errors.New("db error"))
			},
			id: id.String(),
			err: status.Errorf(
				codes.Internal, "cant delete user %v", errors.New("db error")),
		},
	}

	quit := make(chan struct{})
	defer close(quit)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mockdb.NewMockDB(ctrl)

	serv := NewUserService(db, nil)

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
				context.Background(), &pb.DeleteUserRequest{Uuid: v.id})

			assert.ErrorIsf(t, err, v.err, "expected %v, actual: %v", v.err, err)
		})
	}
}
