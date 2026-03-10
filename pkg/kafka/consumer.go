// internal/kafka/consumer.go
package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/pkg/utils/sendto"
)

type OTPConsumerHandler struct{}

type OTPMessage struct {
	OTP   int    `json:"otp"`
	Email string `json:"email"`
}

func (h OTPConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h OTPConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h OTPConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var otpMsg OTPMessage
		if err := json.Unmarshal(msg.Value, &otpMsg); err != nil {
			log.Printf("failed to unmarshal otp message: %v", err)
			session.MarkMessage(msg, "")
			continue
		}
		global.Logger.Info("decode" + otpMsg.Email)
		data := map[string]interface{}{
			"name":       otpMsg.Email,
			"otp":        otpMsg.OTP,
			"expMinutes": 10,
		}

		if err := sendto.SendTemplateEmailOTP([]string{otpMsg.Email}, global.Config.Smtp.User, "verify-email.html", data); err != nil {
			log.Printf("failed to send otp email to %s: %v", otpMsg.Email, err)
		} else {
			log.Printf("OTP email sent to %s", otpMsg.Email)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func StartConsumer(brokers []string, groupID string, topics []string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	defer group.Close()

	ctx := context.Background()
	for {
		if err := group.Consume(ctx, topics, handler); err != nil {
			return err
		}
	}
}
