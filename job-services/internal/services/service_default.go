package services

import (
	"job-services/internal/repositories"
)

type defaultService struct {
	transactionRepo repositories.Repository
}

func NewService(transactionRepo repositories.Repository) Service {
	return &defaultService{
		transactionRepo: transactionRepo,
	}
}