package initialize

import (
	"github.com/kien14502/ecommerce-be/consts"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/pkg/kafka"
)

func InitKafka() {
	brokers := []string{global.Config.Kafka.Host}

	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		global.Logger.Error(err.Error())
		return // ← tránh nil pointer panic
	}

	global.KafkaProducer = producer // ← không defer Close(), producer sống suốt app
	global.Logger.Info("kafka connected!!!")

	go func() {
		if err := kafka.StartConsumer(brokers, "my-group", []string{consts.TopicOTP}, kafka.OTPConsumerHandler{}); err != nil {
			global.Logger.Error(err.Error())
		}
	}()
	// ← bỏ select{}, để main.go quản lý lifecycle
}
