package app

import (
	"job-services/internal/repositories"
	"job-services/internal/services"
	"job-services/internal/worker"
)

func InitPackage() {

	//init repository
	transRepo := repositories.NewRepository(DB.DB().Database(DB_NAME))
	sv := services.NewService(transRepo)

	//init work Group
	workerCount := 5
	batchSize := 10
	workerPool := worker.NewWorkerPool(workerCount, batchSize, sv)

	WorkerPool = workerPool
}
