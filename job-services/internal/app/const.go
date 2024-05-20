package app

import (
	mongo "job-services/mongo"
)

const (
	KAFKA_URL         = "localhost:9091,localhost:9092,localhost:9093"
	KAFKA_TOPIC       = "TEST_TOPIC"
	CONSUMER_GROUP_ID = "TEST_TOPIC_CONSUMER"
	DB_NAME           = "linebk"
)

var MongoConfig = mongo.MongoDBConfig{
	Username:           "admin",
	Password:           "admin",
	Hostname:           "localhost",
	Port:               "27017",
	Options:            "/?authSource=linebk",
	Timeout:            60000,
	MaxPoolSize:        100,
	MinPoolSize:        0,
	MaxIdleConnections: 30000,
}
