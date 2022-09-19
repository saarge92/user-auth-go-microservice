package card

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go-user-microservice/internal/app/card/entities"
	cardErrors "go-user-microservice/internal/app/card/errors"
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

func TestRepositoryCard_ListByCardID(t *testing.T) {
	repositoryCard := getRepositoryTest(t)
	_, e := repositoryCard.ListByCardID(context.Background(), test.UserID)
	require.NoError(t, e)
}

func TestRepositoryCard_OneByCardAndUserID(t *testing.T) {
	repositoryCard := getRepositoryTest(t)
	_, notFoundError := repositoryCard.OneByCardAndUserID(context.Background(), uuid.New().String(), test.UserID)
	require.Error(t, notFoundError)
	require.ErrorIs(t, notFoundError, cardErrors.ErrCardNotFound)
}

func TestRepositoryCard_ExistByCardNumber(t *testing.T) {
	repositoryCard := getRepositoryTest(t)
	notExist, e := repositoryCard.ExistByCardNumber(context.Background(), test.CardNumberForCreate)
	require.NoError(t, e)
	require.False(t, notExist)
}

func getRepositoryTest(t *testing.T) *RepositoryCard {
	require.NoError(t, test.LoadTestEnv())
	dbConnection, closeFunc := test.InitConnectionsWithCloseFunc()
	t.Cleanup(closeFunc)
	databaseWrapper := database.NewDatabase(dbConnection)
	return NewRepositoryCard(databaseWrapper)
}
