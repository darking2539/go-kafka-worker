package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetTransactionByTitle(t *testing.T) {
	// Create a new mock testing instance
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("GetTransaction", func(mt *mtest.T) {

		client := mt.Client
		db := client.Database("test")
		coll := db.Collection("game_transaction")
		repo := &defaultRepository{transaction: coll}

		expectedFilter1 := bson.D{{Key: "title", Value: "transaction1"}}
		expectedFilter2 := bson.D{{Key: "title", Value: "transaction2"}}
		ns := db.Name() + "." + coll.Name()
		resp1 := mtest.CreateCursorResponse(1, ns, mtest.FirstBatch, expectedFilter1)
		resp2 := mtest.CreateCursorResponse(2, ns, mtest.NextBatch, expectedFilter2)
		mt.AddMockResponses(resp1, resp2)

		transactions, err := repo.GetTransactionByTitle("transaction1")
		if err != nil {
			t.Errorf("err: %s", err)
		}

		// Check for errors
		assert.Nil(t, err, "Error fetching transactions: %v", err)

		// Check the number of transactions returned
		assert.Equal(t, 2, len(transactions), "Expected 2 transaction, got %v", len(transactions))
	})
}
