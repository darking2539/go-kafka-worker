package worker

import (
	"encoding/json"
	"kafka-batch/internal/app"
	"kafka-batch/internal/models"
	"log"
	"sync"
)

type WorkerPool struct {
	Workers      []*Worker
	queue        chan *models.TransactionModel
	stop         chan bool
	wg           sync.WaitGroup
	workersCount int
}

type Worker struct {
	queue       chan *models.TransactionModel
	stop        chan bool
	wg          *sync.WaitGroup
	Inserted    int
}

func (p *WorkerPool) Submit(record models.TransactionModel) {
	p.queue <- &record
}

func (p *WorkerPool) Stop() {
	close(p.stop)
	p.wg.Wait()
}

func NewWorkerPool(workersCount int) *WorkerPool {
	pool := &WorkerPool{
		Workers:      make([]*Worker, workersCount),
		queue:        make(chan *models.TransactionModel),
		stop:         make(chan bool),
		workersCount: workersCount,
	}

	pool.wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		worker := NewWorker(pool.queue, pool.stop, &pool.wg)
		worker.Start()
		pool.Workers[i] = worker
	}

	return pool
}

func NewWorker(queue chan *models.TransactionModel, stop chan bool, wg *sync.WaitGroup) *Worker {
	return &Worker{
		queue: queue,
		stop:  stop,
		wg:    wg,
	}
}

func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for {
			select {
			case record := <-w.queue:
				dataBytes, err := json.Marshal(record)
				if err == nil {
					data := models.ProducerRecord{
						Message: dataBytes,
					}
					app.ProducerQueue <- &data
					w.Inserted += 1
				}else {
					log.Println("error: ", err.Error())
				}

			case <-w.stop:
				return
			}
		}
	}()
}
