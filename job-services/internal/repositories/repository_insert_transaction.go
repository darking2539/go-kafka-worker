package repositories

import (
	"job-services/internal/models"
	"job-services/mongo"
)

func (r *defaultRepository) Insert(entity []*models.TransactionModel) error {
	
	ctx, cancelFunc := mongo.InitContext()
	defer cancelFunc()

	// Convert []TransactionModel to []interface{}
	var interfacesSlice []interface{}
	for _, transaction := range entity {
		interfacesSlice = append(interfacesSlice, transaction)
	}

	_, err := r.transaction.InsertMany(ctx, interfacesSlice)
	if err != nil {
		return err
	}

	return nil
}
