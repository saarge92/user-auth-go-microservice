package containers

import (
	"go-user-microservice/internal/forms/builders"
	"go-user-microservice/internal/server"
	"go-user-microservice/internal/services"
	"go.uber.org/dig"
)

func ProvideGrpcServers(container *dig.Container) error {
	e := container.Provide(
		func(userService *services.UserService,
			userFormBuilder *builders.UserFormBuilder,
			jwtService *services.JwtService,
		) *server.UserGrpcServer {
			return server.NewUserGrpcServer(userService, userFormBuilder, jwtService)
		})
	if e != nil {
		return e
	}
	return nil
}
