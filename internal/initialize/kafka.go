package initialize

import "github.com/segmentio/kafka-go"

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

var KafkaProducer *kafka.Writer

func InitKafka() {

}

func CloseKafka() {

}
