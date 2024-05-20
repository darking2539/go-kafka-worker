package worker

import (
	"job-services/internal/models"
	"job-services/internal/services"
	"log"
	"sync"
	"time"
)

func NewWorker(queue chan *models.TransactionModel, stop chan bool, wg *sync.WaitGroup, batchSize int, sv services.Service) *WorkerNode {
	return &WorkerNode{
		service:   sv,
		queue:     queue,
		stop:      stop,
		wg:        wg,
		records:   make([]*models.TransactionModel, 0, batchSize),
		batchSize: batchSize,
	}
}

func (w *WorkerNode) Start() {
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

func (w *WorkerNode) insertedBatch() {

	count := len(w.records)
	err := w.service.AddDataToDB(w.records)

	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Inserted += count       //increment value added
	w.records = w.records[:0] //clear records

}