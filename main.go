package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"user-service/cache"
	"user-service/pb"
	"user-service/repo"
	"user-service/service"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

const (
	postgresConnString = "postgres://root:secret@localhost:5432/users"
)

func main() {
	fmt.Println("server started")

	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", "5000"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pgpool := initPostgres()
	defer pgpool.Close()

	pg := repo.NewPostgres(pgpool)

	redisClient := initRedis()
	defer redisClient.Close()

	redis := cache.NewRedisCache(redisClient)

	userSerivce := service.NewUserService(pg, redis)

	opts := []grpc.ServerOption{
		// grpc.UnaryInterceptor(service.CacheMiddleware),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(grpcServer, userSerivce)
	grpcServer.Serve(lis)

}

func initPostgres() *pgxpool.Pool {
	pgpool, err := pgxpool.Connect(
		context.Background(), postgresConnString)
	if err != nil {
		panic(err)
	}

	return pgpool
}

func initRedis() *redis.Client {
	opt := &redis.Options{
		Network:  "tcp",
		Addr:     net.JoinHostPort("localhost", "6379"),
		DB:       0,
		Password: "secret",
	}

	client := redis.NewClient(opt)
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		panic("cant init redis")
	}

	return client
}
