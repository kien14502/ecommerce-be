package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	kafka "github.com/segmentio/kafka-go"
)

var (
	kafkaProducer *kafka.Writer
)

var (
	host    = "localhost:9092"
	topic   = "example-topic"
	groupID = "example-group"
)

// for producer
func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// for consumer
func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaURL},
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})
}

type StockInfo struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func newStock(msg, typeMsg string) *StockInfo {
	return &StockInfo{
		Message: msg,
		Type:    typeMsg,
	}
}

func actionStock(c *gin.Context) {
	s := newStock(c.Query("msg"), c.Query("type"))
	body := make(map[string]interface{})
	body["action"] = "action"
	body["info"] = s

	json_body, _ := json.Marshal(body)
	msg := kafka.Message{
		Key:   []byte("action"),
		Value: []byte(json_body),
	}
	err := kafkaProducer.WriteMessages(context.Background(), msg)
	if err != nil {
		c.JSON(200, gin.H{
			"Err": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"err": "",
		"msg": "Action Successfully",
	})
}

func RegisterConsumerATC(id int) {
	kafkaGroupId := "consumer-group-"
	reader := getKafkaReader(host, topic, kafkaGroupId)
	defer reader.Close()

	fmt.Printf("Consumer %d started\n", id)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Consumer %d error: %v\n", id, err)
			continue
		}
		fmt.Printf("Consumer %d received message: %s\n", id, string(m.Value))
	}
}

func Main() {
	r := gin.Default()
	kafkaProducer = getKafkaWriter(host, topic)
	defer kafkaProducer.Close()

	r.GET("/kafka/stock", actionStock)

	go RegisterConsumerATC(1)
	go RegisterConsumerATC(2)

	r.Run(":8080")
}
