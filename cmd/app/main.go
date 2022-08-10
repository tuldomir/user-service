package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"user-service/internal/repo"
	"user-service/internal/service"
	"user-service/pb"
	"user-service/pkg/broker"
	"user-service/pkg/cache"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("server started")

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	// host := "localhost"
	// port := "5555"

	lis, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	pgpool := initPostgres()
	defer pgpool.Close()

	pg := repo.NewPostgres(pgpool)

	redisClient := initRedis()
	defer redisClient.Close()
	redisCache := cache.NewRedisCache(redisClient)

	userSerivce := service.NewUserService(pg, redisCache)

	producer := initKafkaProducer()
	kafkaBroker := broker.NewKafkaBroker(producer)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			service.CacheMiddleware(redisCache),
			service.KafkaMiddleware(kafkaBroker),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(grpcServer, userSerivce)
	reflection.Register(grpcServer)
	log.Fatalf("server stopped with err: %v\n", grpcServer.Serve(lis))
}
