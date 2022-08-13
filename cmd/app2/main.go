package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
)

func main() {
	c := initKafkaProducer()
	c.Close()
}

func initRedis() *redis.Client {
	log.Println("init redis")
	// host := os.Getenv("REDIS_HOST")
	host := "localhost"
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

func initKafkaProducer() sarama.SyncProducer {
	log.Println("init kafka")
	brokerCfg := sarama.NewConfig()
	brokerCfg.Producer.RequiredAcks = sarama.WaitForAll
	brokerCfg.Producer.Return.Successes = true

	// host := os.Getenv("KAFKA_HOST")
	host := "localhost"
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
