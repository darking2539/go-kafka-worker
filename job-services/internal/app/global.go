package app

import (
	kf "job-services/kafka"
	"job-services/mongo"

	"github.com/segmentio/kafka-go"
)

var (
	DB *mongo.MongoClient
	ConsumerRunning bool
	ConsumerMsgChan chan *kafka.Message
	ConsumerReader *kf.Cousumer
)

func init() {
	ConsumerRunning = true
	ConsumerMsgChan = make(chan *kafka.Message)
	ConsumerReader = kf.GetConsumerClient(KAFKA_TOPIC, KAFKA_URL, CONSUMER_GROUP_ID)
	InitMongoDBClient(MongoConfig)
}
