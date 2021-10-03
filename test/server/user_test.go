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
	server := test.CreateTestServer()
	container := server.GetDIContainer()
	var userGrpcServer *grpcServer.UserGrpcServer
	e := container.Invoke(
		func(userServer *grpcServer.UserGrpcServer) {
			userGrpcServer = userServer
		})
	assert.Nil(t, e)
	message := &user.SignUpMessage{
		Login:    test.Login,
		Inn:      test.Inn,
		Password: faker.Password(),
		Name:     faker.Name(),
	}
	emptyContext := context.Background()
	response, e := userGrpcServer.Signup(emptyContext, message)
	assert.Nil(t, e)
	assert.NotNil(t, response)
	assert.IsType(t, int32(0), response.Id)
	assert.IsType(t, string(""), response.Token)
}
