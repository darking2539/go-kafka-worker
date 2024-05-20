package worker

import (
	"job-services/internal/models"
	"job-services/internal/services"
	"sync"
)

type WorkerPool struct {
	Workers      []*WorkerNode
	queue        chan *models.TransactionModel
	stop         chan bool
	wg           sync.WaitGroup
	workersCount int
	batchSize    int
}

type WorkerNode struct {
	service   services.Service
	queue     chan *models.TransactionModel
	stop      chan bool
	wg        *sync.WaitGroup
	records   []*models.TransactionModel
	batchSize int
	Inserted  int
}