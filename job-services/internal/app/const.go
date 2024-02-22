package app

import (
	mongo "job-services/mongo"
)

const (
	KAFKA_URL         = "localhost:9091,localhost:9092,localhost:9093"
	KAFKA_TOPIC       = "run-batch-save-db"
	CONSUMER_GROUP_ID = "consumer-job-services"
	DB_NAME           = "hello"
)

var MongoConfig = mongo.MongoDBConfig{
	Username:           "admin",
	Password:           "admin",
	Hostname:           "localhost",
	Port:               "27017",
	Options:            "/?authSource=hello",
	Timeout:            60000,
	MaxPoolSize:        100,
	MinPoolSize:        0,
	MaxIdleConnections: 30000,
}
