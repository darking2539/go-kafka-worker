package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func GetProducer(topic string, kafkaUrl string) *Producer{
	address := strings.Split(kafkaUrl, ",")
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(address...),
			Topic:    topic,
			Balancer: kafka.CRC32Balancer{},
		},
	}
}

func (producer *Producer) PublishMessage(msg []byte) error {

	err := producer.Writer.WriteMessages(context.Background(), kafka.Message{Value: msg})
	if err != nil {
		return err
	}

	return nil
}