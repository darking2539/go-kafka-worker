package services

import "job-services/internal/models"

type Service interface {
	AddDataToDB(transactions []*models.TransactionModel) error
}