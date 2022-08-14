package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	"user-service/internal/mapper"
	"user-service/internal/pb"
	mockdb "user-service/internal/repo/mock"
	"user-service/models"
	"user-service/pkg/broker"
	"user-service/pkg/cache"

	saramamock "github.com/Shopify/sarama/mocks"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCacheMiddleware(t *testing.T) {
	var (
		client      pb.UserServiceClient
		redisClient *redis.Client
	)

	uid := uuid.NewString()
	users := []*models.User{
		{
			UID:       uid,
			Email:     "some",
			CreatedAt: time.Now().UTC(),
		},
	}

	testCases := []struct {
		name        string
		makeStubs   func(*mockdb.MockDB)
		checkResult func()
	}{
		{
			name: "emptyChache",
			makeStubs: func(db *mockdb.MockDB) {

				db.EXPECT().
					List(gomock.Any()).
					Times(1).Return(users, nil)
			},
			checkResult: func() {
				str, err := redisClient.Get(context.Background(), userListKey).Result()
				assert.Error(t, err)
				assert.Empty(t, str)

				res, err := client.List(context.Background(), &emptypb.Empty{})

				assert.NoError(t, err)
				assert.NotEmpty(t, res)

				bs, err := redisClient.Get(context.Background(), userListKey).Bytes()
				assert.NoError(t, err)
				assert.NotEmpty(t, bs)

				var expected []*models.User
				err = json.Unmarshal(bs, &expected)
				assert.NoError(t, err)

				assert.Equal(t, users, expected)
				err = redisClient.Del(context.Background(), userListKey).Err()
				assert.NoError(t, err)
			},
		},
		{
			name: "notEmptyChache",
			makeStubs: func(db *mockdb.MockDB) {

				db.EXPECT().
					List(gomock.Any()).
					Times(0)
			},
			checkResult: func() {
				bs, err := json.Marshal(users)
				assert.NoError(t, err)

				err = redisClient.Set(
					context.Background(), userListKey, bs, 2*time.Second).Err()
				assert.NoError(t, err)

				res, err := client.List(context.Background(), &emptypb.Empty{})
				assert.NoError(t, err)
				assert.NotEmpty(t, res)

				expected, err := mapper.ProtoToUserList(res.Users)
				assert.NoError(t, err)
				assert.Equal(t, users, expected)
				err = redisClient.Del(context.Background(), userListKey).Err()
				assert.NoError(t, err)
			},
		},
		{
			name: "dbChanged/ClearChache/notListMeth",
			makeStubs: func(db *mockdb.MockDB) {

				db.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			checkResult: func() {
				bs, err := json.Marshal(users)
				assert.NoError(t, err)

				err = redisClient.Set(
					context.Background(), userListKey, bs, 5*time.Second).Err()
				assert.NoError(t, err)

				_, err = client.Delete(
					context.Background(), &pb.DeleteUserRequest{Uuid: uuid.NewString()})
				assert.NoError(t, err)

				str, err := redisClient.Get(context.Background(), userListKey).Result()
				assert.Error(t, err)
				assert.Empty(t, str)

				err = redisClient.Del(context.Background(), userListKey).Err()
				assert.NoError(t, err)
			},
		},
	}

	redisServer, err := miniredis.Run()
	assert.NoError(t, err)

	redisClient = redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	defer redisClient.Close()
	redisCache := cache.NewRedisCache(redisClient)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mockdb.NewMockDB(ctrl)

	serv := NewUserService(db)

	quit := make(chan struct{})
	midlew := CacheMiddleware(redisCache)
	defer close(quit)
	bufDialer := bufListener(serv, quit, grpc.UnaryInterceptor(midlew))

	conn, err := grpc.DialContext(
		context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure())

	assert.NoErrorf(t, err, "failed to dial bufnet %v", err)
	defer conn.Close()

	client = pb.NewUserServiceClient(conn)

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.makeStubs(db)

			tC.checkResult()
		})
	}
}

func TestKafkaMiddleware(t *testing.T) {

	user := &models.User{
		UID:       uuid.NewString(),
		Email:     "bob@mail",
		CreatedAt: time.Now().UTC(),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mockdb.NewMockDB(ctrl)

	serv := NewUserService(db)

	kafkaProd := saramamock.NewSyncProducer(t, nil).
		ExpectSendMessageWithCheckerFunctionAndSucceed((func(val []byte) error {

			msg := models.UserEvent{
				EventType: userAddedEvent,
				UID:       user.UID,
				CreatedAt: user.CreatedAt,
			}

			bs, err := json.Marshal(msg)
			if err != nil {
				return err
			}

			if string(bs) != string(val) {
				return errors.New("not equal message values")
			}

			return nil
		}))

	defer func() {
		if err := kafkaProd.Close(); err != nil {
			assert.NoError(t, err)
		}
	}()

	br := broker.NewKafkaBroker(kafkaProd)
	middlw := KafkaMiddleware(br)

	quit := make(chan struct{})
	defer close(quit)
	bufDialer := bufListener(serv, quit, grpc.UnaryInterceptor(middlw))

	conn, err := grpc.DialContext(
		context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure())

	assert.NoErrorf(t, err, "failed to dial bufnet %v", err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	db.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Times(1).Return(user, nil)

	_, err = client.
		Add(context.Background(), &pb.AddUserRequest{Email: user.Email})

	assert.NoError(t, err)
}
