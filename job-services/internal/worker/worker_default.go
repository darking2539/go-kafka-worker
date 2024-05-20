package worker

import (
	"job-services/internal/models"
	"job-services/internal/services"
)

type defaultWorker struct {
	service services.Service
	worker  *WorkerPool
}

func NewWorkerPool(workersCount int, batchSize int, sv services.Service) Worker {
	pool := &WorkerPool{
		Workers:      make([]*WorkerNode, workersCount),
		queue:        make(chan *models.TransactionModel),
		stop:         make(chan bool),
		workersCount: workersCount,
		batchSize:    batchSize,
	}

	pool.wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		worker := NewWorker(pool.queue, pool.stop, &pool.wg, batchSize, sv)
		worker.Start()
		pool.Workers[i] = worker
	}

	return &defaultWorker{
		service: sv,
		worker:  pool,
	}
}
