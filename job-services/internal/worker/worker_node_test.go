package worker

import (
	"job-services/internal/models"
	"job-services/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkerNode_InsertedBatch(t *testing.T) {
	// Create a mock service
	service := mocks.NewService(t)
	service.On("AddDataToDB", mock.Anything).Return(nil)

	// Create a worker node with the mock service
	worker := &WorkerNode{
		service: service,
		records: []*models.TransactionModel{
			{
				Title:       "Test Game",
				Genre:       "Test Genre",
				ReleaseYear: "2022",
				Platform:    "Test Platform",
				Developer:   "Test Developer",
				Publisher:   "Test Publisher",
				Rating:      "Test Rating",
			},
		},
	}

	// Call the insertedBatch method
	worker.insertedBatch()

	// Assert that the AddDataToDB method was called with the correct records
	service.AssertCalled(t, "AddDataToDB", []*models.TransactionModel{
		{
			Title:       "Test Game",
			Genre:       "Test Genre",
			ReleaseYear: "2022",
			Platform:    "Test Platform",
			Developer:   "Test Developer",
			Publisher:   "Test Publisher",
			Rating:      "Test Rating",
		},
	})

	// Assert that the Inserted field was incremented
	assert.Equal(t, 1, worker.Inserted)

	// Assert that the records slice was cleared
	assert.Empty(t, worker.records)
}
