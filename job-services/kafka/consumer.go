package kafka

import (
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Cousumer struct {
	Reader *kafka.Reader
}

func GetConsumerClient(topic string, kafkaUrl string, consumerGroupId string) *Cousumer {
	address := strings.Split(kafkaUrl, ",")
	return &Cousumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        address,
			GroupID:        consumerGroupId,
			Topic:          topic,
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second, // flushes commits to Kafka every second
		}),
	}
}