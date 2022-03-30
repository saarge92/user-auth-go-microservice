package integration

import (
	"context"
	"github.com/bxcodec/faker/v3"
	_ "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/pkg/protobuf/core"
	"go-user-microservice/test"
	"go-user-microservice/test/mocks/providers"
	"go-user-microservice/test/mocks/services"
	"testing"
)

func TestUserSignInSignUp(t *testing.T) {
	stripeServiceProvider := &providers.TestStripeServiceProvider{
		CardStripeServiceMock:    nil,
		AccountStripeServiceMock: &services.AccountStripeServiceMock{},
	}
	serverProvider, closeFunc := test.CreateTestServer(stripeServiceProvider)
	defer closeFunc()
	userGrpcServer := serverProvider.UserGrpcServer()

	password := faker.Password()
	var token string
	t.Run("Should return success sign up messages", func(t *testing.T) {
		message := &core.SignUpMessage{
			Login:    test.LoginForTest,
			Inn:      test.InnForTest,
			Password: password,
			Name:     faker.Name(),
			Country:  "RU",
		}
		emptyContext := context.Background()
		response, e := userGrpcServer.Signup(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &core.SignUpResponse{}, response)
		assert.IsType(t, uint64(0), response.Id)
		assert.IsType(t, "", response.Token)
	})

	t.Run("Should return sign in message", func(t *testing.T) {
		message := &core.SignInMessage{
			Login:    test.LoginForTest,
			Password: password,
		}
		emptyContext := context.Background()
		response, e := userGrpcServer.SignIn(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &core.SignInResponse{}, response)
		token = response.Token
	})

	t.Run("Should verify token properly", func(t *testing.T) {
		message := &core.VerifyMessage{Token: token}
		emptyContext := context.Background()
		response, e := userGrpcServer.VerifyToken(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &core.VerifyMessageResponse{}, response)
		assert.Equal(t, response.User.Login, test.LoginForTest)
	})
}
