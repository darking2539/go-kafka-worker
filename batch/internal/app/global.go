package app

import (
	"kafka-batch/internal/models"
	"kafka-batch/kafka"
)

var (
	ProducerQueue  chan *models.ProducerRecord
	ProducerWriter *kafka.Producer
)

func init() {
	ProducerQueue = make(chan *models.ProducerRecord)
	ProducerWriter = kafka.GetProducer(KAFKA_TOPIC, KAFKA_URL)
}
