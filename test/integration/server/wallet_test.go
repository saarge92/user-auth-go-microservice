package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/internal/app/middlewares/dictionary"
	"go-user-microservice/pkg/protobuf/wallet"
	"go-user-microservice/test"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	serverTest, closeFunc, e := test.CreateTestServer()
	defer closeFunc()
	assert.Nil(t, e)
	walletServer, e := serverTest.GetWalletGrpcServer()
	assert.Nil(t, e)

	t.Run("Should Successfully create", func(t *testing.T) {
		message := &wallet.CreateWalletMessage{
			Code: test.CurrencyCode,
		}
		requestContext := context.WithValue(context.Background(), dictionary.UserID, test.UserIDForRealUser)
		response, e := walletServer.CreateWallet(requestContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
	})
}
