package integration

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/pkg/protobuf/payment"
	"go-user-microservice/test"
	"go-user-microservice/test/mocks/providers"
	"go-user-microservice/test/mocks/services"
	"testing"
)

func TestDeposit(t *testing.T) {
	stripeServiceProvider := &providers.TestStripeServiceProvider{
		CardChargeServiceMock: &services.StripeChargeServiceMock{},
	}
	serverTest, closeFunc := test.CreateTestServer(stripeServiceProvider)
	defer closeFunc()
	paymentGrpcServer := serverTest.PaymentGrpcServer()

	user := &entities.User{ID: test.UserIDForRealUser}
	ctx := context.WithValue(context.Background(), dictionary.User, user)

	t.Run("Deposit wallet from card", func(t *testing.T) {
		paymentDepositRequest := &payment.DepositRequest{
			Amount:           1,
			CardExternalId:   test.ExternalIDForCard,
			WalletExternalId: test.ExternalIDForWallet,
		}
		response, e := paymentGrpcServer.Deposit(ctx, paymentDepositRequest)

		assert.Nil(t, e)
		assert.IsType(t, &payment.DepositResponse{}, response)

		_, uuidParseError := uuid.Parse(response.TransactionId)
		assert.Equal(t, true, uuidParseError == nil, "Response is not uuid type")
	})
}
