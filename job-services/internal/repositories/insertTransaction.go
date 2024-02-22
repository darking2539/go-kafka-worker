package repositories

import (
	"job-services/internal/app"
	"job-services/internal/models"
	"job-services/mongo"
)

func Insert(entity []*models.TransactionModel) error {
	
	ctx, cancelFunc := mongo.InitContext()
	defer cancelFunc()

	// Convert []TransactionModel to []interface{}
	var interfacesSlice []interface{}
	for _, transaction := range entity {
		interfacesSlice = append(interfacesSlice, transaction)
	}

	_, err := app.DB.DB().Database(app.DB_NAME).Collection("game_transaction").InsertMany(ctx, interfacesSlice)
	if err != nil {
		return err
	}

	return nil
}
