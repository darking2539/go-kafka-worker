package worker

import (
	"job-services/internal/models"
	"job-services/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubmit(t *testing.T) {
	service := mocks.NewService(t)
	service.On("AddDataToDB", mock.Anything).Return(nil)
	worker := NewWorkerPool(1, 1, service)
	record := &models.TransactionModel{
		Title:       "Test Game",
		Genre:       "Test Genre",
		ReleaseYear: "2022",
		Platform:    "Test Platform",
		Developer:   "Test Developer",
		Publisher:   "Test Publisher",
		Rating:      "Test Rating",
	}

	worker.Submit(record)
	close(worker.GetWorkerPool().stop)
	worker.GetWorkerPool().wg.Wait()

	service.AssertCalled(t, "AddDataToDB", []*models.TransactionModel{
		record,
	})

}

func TestStop(t *testing.T) {

	sv := mocks.NewService(t)
	worker := NewWorkerPool(1, 1, sv)

	worker.Stop()

	// Assert that the worker pool has finished processing all tasks
	worker.GetWorkerPool().wg.Wait()

	_, ok := <-worker.GetWorkerPool().queue
	assert.False(t, ok, "Channel is not close then error.")

	_, ok = <-worker.GetWorkerPool().stop
	assert.False(t, ok, "Channel is not close then error.")

}
