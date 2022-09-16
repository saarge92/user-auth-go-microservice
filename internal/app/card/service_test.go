package card

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/card"
	"github.com/stripe/stripe-go/v72/client"
	"github.com/stripe/stripe-go/v72/token"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/card/mocks"
	userDto "go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/internal/pkg/dictionary"
	stripeServices "go-user-microservice/internal/pkg/services/stripe"
	"go-user-microservice/pkg/protobuf/core"
	"go-user-microservice/test"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type testCardServiceStruct struct {
	serviceCard   *ServiceCard
	stripeBackend *mocks.StripeBackend
}

func TestServiceCard_Create(t *testing.T) {
	testStructData := getServiceTestStruct(t)
	serviceCard := testStructData.serviceCard
	stripeBackend := testStructData.stripeBackend

	createCard := forms.CreateCard{CreateCardRequest: &core.CreateCardRequest{
		CardNumber:  test.CardNumberForCreate,
		ExpireMonth: 03,
		ExpireYear:  uint32(time.Now().Year() + 2),
		Cvc:         333,
		IsDefault:   true,
	}}

	userRoleDTO := &userDto.UserRole{
		User: entities.User{
			ID:                 test.UserID,
			AccountProviderID:  test.UserAccountProviderID,
			CustomerProviderID: test.UserCustomerID,
		},
	}
	ctx := context.WithValue(context.Background(), dictionary.User, userRoleDTO)

	expireMonth := strconv.Itoa(int(createCard.ExpireMonth))
	expireYear := strconv.Itoa(int(createCard.ExpireYear))
	cvc := strconv.Itoa(int(createCard.Cvc))
	tokenParamsExpected := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Account:  stripe.String(userRoleDTO.User.AccountProviderID),
			Number:   stripe.String(createCard.CardNumber),
			ExpMonth: stripe.String(expireMonth),
			ExpYear:  stripe.String(expireYear),
			CVC:      stripe.String(cvc),
			Currency: stripe.String("USD"),
		},
	}

	stripeBackend.EXPECT().Call(http.MethodPost, "/v1/tokens", mock.Anything, tokenParamsExpected, mock.Anything).
		Return(nil)

	cardExpected := &stripe.Card{ID: uuid.New().String()}
	stripeBackend.On("CallRaw", http.MethodPost, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.MatchedBy(func(cardElement interface{}) bool {
			cardDataPassed := cardElement.(*stripe.Card)
			*cardDataPassed = *cardExpected
			return true
		})).Return(nil)

	response, e := serviceCard.Create(ctx, createCard)
	require.NoError(t, e)
	require.NotEmpty(t, response)
}

func getServiceTestStruct(t *testing.T) testCardServiceStruct {
	require.NoError(t, test.LoadTestEnv())
	dbConnection, closeFunc := test.InitConnectionsWithCloseFunc()
	t.Cleanup(closeFunc)

	databaseWrapper := database.NewDatabase(dbConnection)
	cardRepository := NewRepositoryCard(databaseWrapper)

	stripeClient := &client.API{}

	stripeBackend := &mocks.StripeBackend{}
	stripeClient.Tokens = &token.Client{B: stripeBackend}
	stripeClient.Cards = &card.Client{B: stripeBackend}

	stripeService := stripeServices.NewCardStripeService(stripeClient)
	serviceCard := NewServiceCard(cardRepository, stripeService)

	return testCardServiceStruct{
		serviceCard:   serviceCard,
		stripeBackend: stripeBackend,
	}
}
