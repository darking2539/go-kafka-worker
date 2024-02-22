package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	mgc    *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func ConnectMongoDBByConfig(cfg MongoDBConfig) (*MongoClient, error) {
	timeout := int64(10000)
	if cfg.Timeout != 0 {
		timeout = cfg.Timeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)

	URI := fmt.Sprintf("mongodb://%s:%s@%s:%s%s", cfg.Username, cfg.Password, cfg.Hostname, cfg.Port, cfg.Options)

	opts := options.Client().ApplyURI(URI).SetMaxConnIdleTime(time.Duration(cfg.MaxIdleConnections) * time.Millisecond).SetMaxPoolSize(uint64(cfg.MaxPoolSize)).SetMinPoolSize(uint64(cfg.MinPoolSize))
	conn, err := mongo.Connect(ctx, opts)
	if err != nil {
		cancel()
		return nil, err
	}

	client := MongoClient{
		mgc:    conn,
		ctx:    ctx,
		cancel: cancel,
	}

	return &client, nil
}

func (c *MongoClient) Close() {
	err := c.mgc.Disconnect(c.ctx)
	if err != nil {
		log.Panicf("cfg.DBConn.Disconnect: %v", err)
	}
	c.cancel()
}

func (c *MongoClient) DB() *mongo.Client {
	return c.mgc
}

func (c *MongoClient) Context() context.Context {
	return c.ctx
}

func (c *MongoClient) Ping(rp *readpref.ReadPref) error {
	err := c.mgc.Ping(c.ctx, rp)
	if err != nil {
		return err
	}
	return nil
}
