package card

import (
	"context"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v72"
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

func TestServiceCard_Create(t *testing.T) {
	testData := getServerTestStruct(t)
	serverInstance := testData.server
	stripeServiceCard := testData.stripeServiceCard

	ctx := context.WithValue(context.Background(), dictionary.User, test.UserRoleData)

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
		AccountProviderID: test.UserRoleData.User.AccountProviderID,
	}
	cardResponseExpected := &stripe.Card{
		ID: uuid.New().String(),
	}
	stripeServiceCard.EXPECT().CreateCard(cardParameterExpected).Return(cardResponseExpected, nil)

	cardResponse, e := serverInstance.CreateCard(ctx, createCard)
	require.NoError(t, e)
	require.NoError(t, is.UUID.Validate(cardResponse.ExternalId), "card's external_id is not uuid")
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
