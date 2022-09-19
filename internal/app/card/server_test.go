package card

import (
	"context"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v72"
	cardErrors "go-user-microservice/internal/app/card/errors"
	"go-user-microservice/internal/app/card/mocks"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/pkg/protobuf/core"
	"go-user-microservice/test"
	"testing"
	"time"
)

type serverTestStruct struct {
	server            *GrpcServerCard
	stripeServiceCard *mocks.StripeCardService
}

func TestServiceCard_Create_List(t *testing.T) {
	testData := getServerTestStruct(t)
	serverInstance := testData.server
	stripeServiceCard := testData.stripeServiceCard

	currentUser := test.CurrentUser
	ctx := context.WithValue(context.Background(), dictionary.CurrentUser, currentUser)

	t.Run("Create should be success", func(t *testing.T) {
		createCard := &core.CreateCardRequest{
			CardNumber:  test.CardNumberForCreate,
			ExpireMonth: 03,
			ExpireYear:  uint32(time.Now().Year() + 2),
			Cvc:         333,
			IsDefault:   true,
		}

		cardParameterExpected := dto.StripeCardCreate{
			Number:            createCard.CardNumber,
			ExpireMonth:       uint8(createCard.ExpireMonth),
			ExpireYear:        createCard.ExpireYear,
			CVC:               createCard.Cvc,
			AccountProviderID: currentUser.AccountProviderID,
		}
		cardResponseExpected := &stripe.Card{
			ID: uuid.New().String(),
		}
		stripeServiceCard.EXPECT().CreateCard(cardParameterExpected).Return(cardResponseExpected, nil)

		cardResponse, e := serverInstance.CreateCard(ctx, createCard)
		require.NoError(t, e)
		require.NoError(t, is.UUID.Validate(cardResponse.ExternalId), "card's external_id is not uuid")
	})

	t.Run("MyCards should be success", func(t *testing.T) {
		cards, e := serverInstance.MyCards(ctx, &empty.Empty{})
		require.NoError(t, e)
		require.True(t, len(cards.Cards) > 0, "Cards list is empty")
	})

	t.Run("CreateCard should return already exist", func(t *testing.T) {
		createCard := &core.CreateCardRequest{
			CardNumber:  test.CardNumberForCreate,
			ExpireMonth: 03,
			ExpireYear:  uint32(time.Now().Year() + 2),
			Cvc:         333,
			IsDefault:   true,
		}

		_, e := serverInstance.CreateCard(ctx, createCard)
		require.ErrorIs(t, e, cardErrors.ErrCardAlreadyExist)
	})
}

func getServerTestStruct(t *testing.T) serverTestStruct {
	require.NoError(t, test.LoadTestEnv())
	dbConnection, closeFunc := test.InitConnectionsWithCloseFunc()
	t.Cleanup(closeFunc)

	databaseWrapper := database.NewDatabase(dbConnection)
	transactionHandler := database.NewTransactionHandler(dbConnection)
	cardRepository := NewRepositoryCard(databaseWrapper)
	stripeServiceCard := &mocks.StripeCardService{}
	serviceCard := NewServiceCard(cardRepository, stripeServiceCard)

	serverInstance := NewGrpcServerCard(serviceCard, transactionHandler)

	return serverTestStruct{
		server:            serverInstance,
		stripeServiceCard: stripeServiceCard,
	}
}
