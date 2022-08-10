package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

func initKafkaProducer() sarama.SyncProducer {
	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	host := os.Getenv("KAFKA_HOST")
	// host := "localhost"
	port := os.Getenv("KAFKA_PORT")

	addr := net.JoinHostPort(host, port)
	producer, err := sarama.NewSyncProducer(
		[]string{
			addr,
		}, brokerCfg)

	if err != nil {
		log.Fatalf("cant init kafka producer %v", err)
	}

	return producer

}

func initPostgres() *pgxpool.Pool {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	// host := "localhost"
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v", user, pass, host, port, db)

	pgpool, err := pgxpool.Connect(
		context.Background(), connStr)
	if err != nil {
		log.Fatalf("cant init postgres %v", err)
	}

	return pgpool
}

func initRedis() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	// host := "localhost"
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASSWORD")
	strDB := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(strDB)
	if err != nil {
		log.Fatalf("cant init redis %v\n", err)
	}

	addr := net.JoinHostPort(host, port)
	opt := &redis.Options{
		Network:  "tcp",
		Addr:     addr,
		DB:       db,
		Password: pass,
	}

	client := redis.NewClient(opt)
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		log.Fatalf("cant init redis %v\n", err)
	}

	return client
}
