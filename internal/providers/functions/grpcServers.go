package functions

import (
	"go-user-microservice/internal/server"
	"go-user-microservice/internal/services"
	"go.uber.org/dig"
)

func ProvideGrpcServers(container *dig.Container) error {
	e := container.Provide(func(userService *services.UserService) *server.UserGrpcServer {
		return server.NewUserGrpcServer(userService)
	})
	if e != nil {
		return e
	}
	return nil
}
