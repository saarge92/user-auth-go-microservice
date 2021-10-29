package server

import (
	"context"
	"github.com/bxcodec/faker/v3"
	_ "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"go-user-microservice/pkg/protobuf/user_server"
	"go-user-microservice/test"
	"testing"
)

func TestUserSignInSignUp(t *testing.T) {
	server, closeFunc, e := test.CreateTestServer()
	assert.Nil(t, e)
	defer closeFunc()
	userGrpcServer, e := server.GetUserGrpcServer()
	assert.Nil(t, e)
	password := faker.Password()
	var token string
	t.Run("Should return success sign up messages", func(t *testing.T) {
		message := &user_server.SignUpMessage{
			Login:    test.Login,
			Inn:      test.Inn,
			Password: password,
			Name:     faker.Name(),
		}
		emptyContext := context.Background()
		response, e := userGrpcServer.Signup(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &user_server.SignUpResponse{}, response)
		assert.IsType(t, uint64(0), response.Id)
		assert.IsType(t, string(""), response.Token)
	})

	t.Run("Should return sign in message", func(t *testing.T) {
		message := &user_server.SignInMessage{
			Login:    test.Login,
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
		assert.Equal(t, response.User.Login, test.Login)
	})
}
