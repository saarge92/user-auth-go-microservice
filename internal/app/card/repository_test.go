package card

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/test"
	"testing"
	"time"
)

func TestRepositoryCard_Create(t *testing.T) {
	repositoryCard := getRepositoryTest(t)
	cardEntity := &entities.Card{
		UserID:             test.UserID,
		Number:             test.CardNumberForCreate,
		ExternalProviderID: uuid.New().String(),
		ExpireYear:         uint32(time.Now().Year() + 2),
		ExpireMonth:        03,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
	}
	ctx := context.Background()

	e := repositoryCard.Create(ctx, cardEntity)
	require.NoError(t, e)
}

func getRepositoryTest(t *testing.T) *RepositoryCard {
	require.NoError(t, test.LoadTestEnv())
	dbConnection, closeFunc := test.InitConnectionsWithCloseFunc()
	t.Cleanup(closeFunc)
	databaseWrapper := database.NewDatabase(dbConnection)
	return NewRepositoryCard(databaseWrapper)
}
