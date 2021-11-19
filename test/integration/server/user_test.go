package server

import (
	"context"
	"github.com/bxcodec/faker/v3"
	_ "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/pkg/protobuf/user_server"
	"go-user-microservice/test"
	"go-user-microservice/test/test_providers"
	"testing"
)

func TestUserSignInSignUp(t *testing.T) {
	stripeContainerProvider := test_providers.ProvideStripe
	server, closeFunc, e := test.CreateTestServer(stripeContainerProvider)
	assert.Nil(t, e)
	defer closeFunc()
	userGrpcServer := server.GetUserGrpcServer()

	assert.NotNil(t, userGrpcServer)
	password := faker.Password()
	var token string
	t.Run("Should return success sign up messages", func(t *testing.T) {
		message := &user_server.SignUpMessage{
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
		assert.IsType(t, &user_server.SignUpResponse{}, response)
		assert.IsType(t, uint64(0), response.Id)
		assert.IsType(t, "", response.Token)
	})

	t.Run("Should return sign in message", func(t *testing.T) {
		message := &user_server.SignInMessage{
			Login:    test.LoginForTest,
			Password: password,
		}
		emptyContext := context.Background()
		response, e := userGrpcServer.SignIn(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &user_server.SignInResponse{}, response)
		token = response.Token
	})

	t.Run("Should verify token properly", func(t *testing.T) {
		message := &user_server.VerifyMessage{Token: token}
		emptyContext := context.Background()
		response, e := userGrpcServer.VerifyToken(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &user_server.VerifyMessageResponse{}, response)
		assert.Equal(t, response.User.Login, test.LoginForTest)
	})
}
