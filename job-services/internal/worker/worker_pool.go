package worker

import (
	"job-services/internal/models"
)


func (p *defaultWorker) Submit(record *models.TransactionModel) {
	p.worker.queue <- record
}

func (p *defaultWorker) Stop() {
	close(p.worker.stop)
	p.worker.wg.Wait()
	close(p.worker.queue)
}


func (p *defaultWorker) GetWorkerPool() *WorkerPool {
	return p.worker
}