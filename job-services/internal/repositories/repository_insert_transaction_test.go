package repositories

import (
	"job-services/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestInsertTransaction(t *testing.T) {
	// Create a new mock testing instance
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("InsertTransaction", func(mt *mtest.T) {

		client := mt.Client
		db := client.Database("test")
		coll := db.Collection("game_transaction")
		repo := &defaultRepository{transaction: coll}

		resp := mtest.CreateSuccessResponse()
		mt.AddMockResponses(resp)
		
		//mock data
		data := []*models.TransactionModel{{Title:"transaction1"},{Title:"transaction2"}}

		err := repo.Insert(data)
		if err != nil {
			t.Errorf("err: %s", err)
		}

		// Check for errors
		assert.Nil(t, err, "Insert Err: %v", err)

	})
}
