package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/pkg/protobuf/card"
	"go-user-microservice/test"
	"go-user-microservice/test/mocks/providers"
	"go-user-microservice/test/mocks/services"
	"testing"
)

func TestCardAdd(t *testing.T) {
	stripeServiceProvider := &providers.TestServiceProvider{
		CardStripeServiceMock: &services.StripeCardServiceMock{},
	}
	serverTest, closeFunc := test.CreateTestServer(stripeServiceProvider)
	defer closeFunc()
	cardServer := serverTest.CardGrpcServer()

	t.Run("Add Card For User", func(t *testing.T) {
		user := &entities.User{ID: test.UserIDForRealUser}
		ctx := context.WithValue(context.Background(), dictionary.User, user)
		request := &card.CreateCardRequest{
			CardNumber:  test.CardNumber,
			ExpireMonth: test.ExpireMonth,
			ExpireYear:  test.ExpireYear,
			Cvc:         test.CVC,
			IsDefault:   true,
		}
		createCard, e := cardServer.CreateCard(ctx, request)
		assert.Nil(t, e)
		assert.IsType(t, &card.CreateCardResponse{}, createCard)
	})
}
