package kafka

import (
	"context"
	"encoding/json"
	"log"

	kafkago "github.com/segmentio/kafka-go"
)

// MessageHandler is the callback invoked for each Kafka message.
type MessageHandler func(ctx context.Context, event EmailEvent) error

// Consumer wraps a kafka-go reader.
type Consumer struct {
	reader  *kafkago.Reader
	handler MessageHandler
}

// NewConsumer creates a new Kafka consumer for the user.registered topic.
func NewConsumer(brokers []string, groupID string, handler MessageHandler) *Consumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    TopicUserRegistered,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{reader: r, handler: handler}
}

// Start begins consuming messages from Kafka (blocking).
func (c *Consumer) Start(ctx context.Context) {
	log.Println("[Kafka Consumer] Starting consumer for topic:", TopicUserRegistered)
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("[Kafka Consumer] Context cancelled, stopping consumer")
				return
			}
			log.Printf("[Kafka Consumer] Error reading message: %v", err)
			continue
		}

		var event EmailEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("[Kafka Consumer] Error unmarshalling message: %v", err)
			continue
		}

		log.Printf("[Kafka Consumer] Received event for user: %s", event.Email)

		if err := c.handler(ctx, event); err != nil {
			log.Printf("[Kafka Consumer] Error handling event: %v", err)
		}
	}
}

// Close closes the Kafka reader.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
