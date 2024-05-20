package worker

import "job-services/internal/models"

type Worker interface {
	Submit(record *models.TransactionModel)
	Stop()
	GetWorkerPool() *WorkerPool
}
