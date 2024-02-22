package worker

import (
	"job-services/internal/models"
	"job-services/internal/repositories"
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	Workers      []*Worker
	queue        chan *models.TransactionModel
	stop         chan bool
	wg           sync.WaitGroup
	workersCount int
	batchSize    int
}

type Worker struct {
	queue     chan *models.TransactionModel
	stop      chan bool
	wg        *sync.WaitGroup
	records   []*models.TransactionModel
	batchSize int
	Inserted  int
}

func (p *WorkerPool) Submit(record *models.TransactionModel) {
	p.queue <- record
}

func (p *WorkerPool) Stop() {
	close(p.stop)
	p.wg.Wait()
}

func NewWorkerPool(workersCount int, batchSize int) *WorkerPool {
	pool := &WorkerPool{
		Workers:      make([]*Worker, workersCount),
		queue:        make(chan *models.TransactionModel),
		stop:         make(chan bool),
		workersCount: workersCount,
		batchSize:    batchSize,
	}

	pool.wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		worker := NewWorker(pool.queue, pool.stop, &pool.wg, batchSize)
		worker.Start()
		pool.Workers[i] = worker
	}

	return pool
}

func NewWorker(queue chan *models.TransactionModel, stop chan bool, wg *sync.WaitGroup, batchSize int) *Worker {
	return &Worker{
		queue:     queue,
		stop:      stop,
		wg:        wg,
		records:   make([]*models.TransactionModel, 0, batchSize),
		batchSize: batchSize,
	}
}

func (w *Worker) Start() {
	ticker := time.NewTicker(10*time.Second)

	go func() {
		defer w.wg.Done()
		defer ticker.Stop()
		for {
			select {
			case record := <-w.queue:
				w.records = append(w.records, record)
				if len(w.records) >= w.batchSize {
					w.insertedBatch()
				}
			case <-ticker.C:
                if len(w.records) > 0 {
                    w.insertedBatch()
                }
			case <-w.stop:
				if len(w.records) > 0 {
					w.insertedBatch()
				}
				return
			}
		}
	}()
}

func (w *Worker) insertedBatch() {

	count := len(w.records)
	err := repositories.Insert(w.records)

	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Inserted += count       //increment value added
	w.records = w.records[:0] //clear records

}
