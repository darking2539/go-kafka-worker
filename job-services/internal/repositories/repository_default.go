package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type defaultRepository struct {
	transaction *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	gameTransaction := db.Collection("game_transaction")
	return &defaultRepository{
		transaction: gameTransaction,
	}
}
