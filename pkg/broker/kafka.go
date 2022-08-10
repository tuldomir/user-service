package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

// KafkaBroker .
type KafkaBroker struct {
	producer sarama.SyncProducer
}

// NewKafkaBroker .
func NewKafkaBroker(producer sarama.SyncProducer) *KafkaBroker {
	return &KafkaBroker{producer: producer}
}

// Publish .
func (b *KafkaBroker) Publish(
	ctx context.Context, topic string, event interface{}) error {

	bs, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("cant marshal event %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(bs),
	}

	_, _, err = b.producer.SendMessage(msg)

	return err
}
