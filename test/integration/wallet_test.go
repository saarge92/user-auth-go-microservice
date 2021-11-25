package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/pkg/protobuf/wallet"
	"go-user-microservice/test"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	serverTest, closeFunc := test.CreateTestServer(nil)
	defer closeFunc()
	walletServer := serverTest.WalletGrpcServer()

	t.Run("Should Successfully create", func(t *testing.T) {
		message := &wallet.CreateWalletRequest{
			Code: test.CurrencyCode,
		}
		user := &entities.User{
			ID: test.UserIDForRealUser,
		}
		requestContext := context.WithValue(context.Background(), dictionary.User, user)
		response, e := walletServer.CreateWallet(requestContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &wallet.CreateWalletResponse{}, response)
		assert.NotNil(t, response.ExternalId)
	})
}
