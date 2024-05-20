package repositories

import "job-services/internal/models"

type Repository interface {
	Insert(entity []*models.TransactionModel) error
	GetTransactionByTitle(title string) ([]models.TransactionModel, error)
}