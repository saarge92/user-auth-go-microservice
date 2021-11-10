package providers

import (
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/user/forms/builders"
	services2 "go-user-microservice/internal/app/user/services"
	"go.uber.org/dig"
)

func ProvideUserGrpcServers(container *dig.Container) error {
	return container.Provide(
		func(
			userService *services2.ServiceUser,
			authService *services2.AuthService,
		) *user.GrpcUserServer {
			userFormBuilder := &builders.UserFormBuilder{}
			return user.NewUserGrpcServer(userFormBuilder, authService)
		},
	)
}
