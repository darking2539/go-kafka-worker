package app

import (
	mongo "job-services/mongo"
	"log"
)

func InitMongoDBClient(cfg mongo.MongoDBConfig) error {
	//init mongodb
	mgc, err := mongo.ConnectMongoDBByConfig(cfg)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
		return err
	}

	err = mgc.Ping(nil)
	if err != nil {
		log.Fatalf("Cannot ping DB: %v", err)
		return err
	}

	DB = mgc
	return nil
}