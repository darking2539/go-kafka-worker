package repositories

import (
	"job-services/internal/models"
	"job-services/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *defaultRepository) GetTransactionByTitle(title string) ([]models.TransactionModel, error) {
	
	ctx, cancelFunc := mongo.InitContext()
	defer cancelFunc()

	filter := bson.M{"title": title}

	cursor, err := r.transaction.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	transactions := []models.TransactionModel{}

	for cursor.Next(ctx) {
		var elem models.TransactionModel
		err = cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, elem)
	}

	return transactions, nil
}