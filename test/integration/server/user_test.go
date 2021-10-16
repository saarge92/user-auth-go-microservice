package server

import (
	"context"
	"github.com/bxcodec/faker/v3"
	_ "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	server2 "go-user-microservice/internal/app/server"
	"go-user-microservice/pkg/protobuf/user"
	"go-user-microservice/test"
	"testing"
)

func TestUserSignInSignUp(t *testing.T) {
	server, closeFunc, e := test.CreateTestServer()
	assert.Nil(t, e)
	defer closeFunc()
	container := server.GetDIContainer()
	var userGrpcServer *server2.UserGrpcServer
	e = container.Invoke(
		func(userServer *server2.UserGrpcServer) {
			userGrpcServer = userServer
		})
	assert.Nil(t, e)
	password := faker.Password()
	var token string
	t.Run("Should return success sign up messages", func(t *testing.T) {
		message := &user.SignUpMessage{
			Login:    test.Login,
			Inn:      test.Inn,
			Password: password,
			Name:     faker.Name(),
		}
		emptyContext := context.Background()
		response, e := userGrpcServer.Signup(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &user.SignUpResponse{}, response)
		assert.IsType(t, uint64(0), response.Id)
		assert.IsType(t, string(""), response.Token)
	})

	t.Run("Should return sign in message", func(t *testing.T) {
		message := &user.SignInMessage{
			Login:    test.Login,
			Password: password,
		}
		emptyContext := context.Background()
		response, e := userGrpcServer.SignIn(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &user.SignInResponse{}, response)
		token = response.Token
	})

	t.Run("Should verify token properly", func(t *testing.T) {
		message := &user.VerifyMessage{Token: token}
		emptyContext := context.Background()
		response, e := userGrpcServer.VerifyToken(emptyContext, message)
		assert.Nil(t, e)
		assert.NotNil(t, response)
		assert.IsType(t, &user.VerifyMessageResponse{}, response)
		assert.Equal(t, response.User.Login, test.Login)
	})
}
