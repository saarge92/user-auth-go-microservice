package containers

import (
	builders2 "go-user-microservice/internal/app/forms/builders"
	server2 "go-user-microservice/internal/app/server"
	"go-user-microservice/internal/app/services/member"
	"go.uber.org/dig"
)

func ProvideGrpcServers(container *dig.Container) error {
	e := container.Provide(
		func(userService *member.UserService,
			userFormBuilder *builders2.UserFormBuilder,
			authService *member.AuthService,
		) *server2.UserGrpcServer {
			return server2.NewUserGrpcServer(userFormBuilder, authService)
		})
	if e != nil {
		return e
	}
	return nil
}
