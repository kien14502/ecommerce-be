package initialize

import (
	"log"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/segmentio/kafka-go"
)

var kafkaProducer *kafka.Writer

func InitKafka() {
	global.Kafka = &kafka.Writer{
		Addr:     kafka.TCP("localhost:19092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.Kafka.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %v", err)
	}
}
