package server

import (
	"context"
	"github.com/bxcodec/faker/v3"
	_ "github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	grpcServer "go-user-microservice/internal/server"
	"go-user-microservice/pkg/protobuf/user"
	"go-user-microservice/test"
	"testing"
)

func TestUserSignInSignUp(t *testing.T) {
	server, closeFunc, e := test.CreateTestServer()
	assert.Nil(t, e)
	defer closeFunc()
	container := server.GetDIContainer()
	var userGrpcServer *grpcServer.UserGrpcServer
	e = container.Invoke(
		func(userServer *grpcServer.UserGrpcServer) {
			userGrpcServer = userServer
		})
	assert.Nil(t, e)
	password := faker.Password()

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
	})
}
