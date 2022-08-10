package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	type Event struct {
		EventType string    `json:"event_type"`
		CreatedAt time.Time `json:"created_at"`
	}

	var bs []byte
	f := func() {

		// 	uid, err := uuid.NewUUID()
		// 	if err != nil {
		// 		log.Fatalf("Failed gen uui %v", err)
		// 	}

		// 	bs, err = json.Marshal(domain.UserEvent{
		// 		EventType: domain.UserAddedEvent,
		// 		UID:       uid,
		// 		CreatedAt: time.Now().UTC(),
		// 	})
		// 	if err != nil {
		// 		log.Fatalf("Failed to marshal useraddevent %v", err)

		// 	}

		bs, err = json.Marshal(Event{
			EventType: "useradevent",
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Fatalf("Failed to marshal useraddevent %v", err)

		}

	}

	n := 10
	for i := 0; i <= n; i++ {
		f()
		_, _, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "github",
			Value: sarama.StringEncoder(bs),
		})

		if err != nil {
			log.Fatalf("Failed to store your data %v", err)
		}
	}

}
