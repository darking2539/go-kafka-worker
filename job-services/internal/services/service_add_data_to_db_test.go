package services

import (
	"job-services/internal/models"
	"job-services/internal/repositories/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddDataToDB(t *testing.T) {
	// TODO: Write your test code here
	repo := mocks.NewRepository(t)
	repo.On("Insert", mock.Anything).Return(nil) // Add this line to mock the Write function
	
	service := &defaultService{
		transactionRepo: repo,
	}

	game1 := &models.TransactionModel{
        Title:       "The Legend of Zelda: Breath of the Wild",
        Genre:       "Action-Adventure",
        ReleaseYear: "2017",
        Platform:    "Nintendo Switch",
        Developer:   "Nintendo EPD",
        Publisher:   "Nintendo",
        Rating:      "10/10",
    }

    game2 := &models.TransactionModel{
        Title:       "Super Mario Odyssey",
        Genre:       "Platformer",
        ReleaseYear: "2017",
        Platform:    "Nintendo Switch",
        Developer:   "Nintendo EPD",
        Publisher:   "Nintendo",
        Rating:      "9.5/10",
    }

	transactions := []*models.TransactionModel{
		game1, game2,
	}

	err := service.AddDataToDB(transactions)
	if err != nil {
		t.Errorf("err: %s", err)
	}

	assert.Nil(t, err)
}
