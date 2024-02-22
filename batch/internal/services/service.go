package services

import (
	"kafka-batch/internal/models"
	"kafka-batch/internal/worker"
)

func ProcessRecord(record []string, workerPool *worker.WorkerPool, countCSV *int) error {

	transactionData := models.TransactionModel{
		Title:       record[0],
		Genre:       record[1],
		ReleaseYear: record[2],
		Platform:    record[3],
		Developer:   record[4],
		Publisher:   record[5],
		Rating:      record[6],
	}

	workerPool.Submit(transactionData)
	*countCSV++

	return nil
}
