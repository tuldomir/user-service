package broker

import "context"

// Broker .
type Broker interface {
	Publish(ctx context.Context, topic string, event interface{}) error
}
