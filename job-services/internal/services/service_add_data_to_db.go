package services

import "job-services/internal/models"

func (s *defaultService) AddDataToDB(transactions []*models.TransactionModel) error {

	err := s.transactionRepo.Insert(transactions)
	if err != nil {
		return err
	}

	return nil
}
