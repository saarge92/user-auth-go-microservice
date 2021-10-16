package containers

import (
	"go-user-microservice/internal/app/middlewares"
	"go-user-microservice/internal/app/services/member"
	"go.uber.org/dig"
)

func ProvideUserMiddlewares(container *dig.Container) error {
	e := container.Provide(
		func(authService *member.AuthService) *middlewares.UserGrpcMiddleware {
			return middlewares.NewUserGrpcMiddleware(authService)
		},
	)
	return e
}
