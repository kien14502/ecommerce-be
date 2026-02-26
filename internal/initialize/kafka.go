package initialize

import (
	"log"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/pkg/kafka"
)

const (
	TopicOTP               = "otp-auth-topic"
	TopicEmailNotification = "email-notification-topic"
	TopicOrderCreated      = "order-created-topic"
	TopicOrderPaid         = "order-paid-topic"
	TopicOrderShipped      = "order-shipped-topic"
	TopicOrderCancelled    = "order-cancelled-topic"
	TopicPaymentSuccess    = "payment-success-topic"
	TopicPaymentFailed     = "payment-failed-topic"
	TopicInventoryUpdated  = "inventory-updated-topic"
	TopicUserRegistered    = "user-registered-topic"
)

func InitKafka() {
	kafkaHost := global.Config.Kafka.Host
	brokers := []string{kafkaHost}
	global.KafkaManager = kafka.NewKafkaManager(brokers)
	log.Println("Kafka manager initialized")
}

func CloseKafka() {
	if err := global.KafkaManager.Close(); err != nil {
		log.Printf("Failed to close kafka manager: %v", err)
	}
}
