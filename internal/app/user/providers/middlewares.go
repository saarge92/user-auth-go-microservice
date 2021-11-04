package providers

import (
	"go-user-microservice/internal/app/user/middlewares"
	"go-user-microservice/internal/app/user/services"
	"go.uber.org/dig"
)

func ProvideGrpcMiddleware(container *dig.Container) error {
	e := container.Provide(
		func(jwtService *services.JwtService) *middlewares.UserGrpcMiddleware {
			return middlewares.NewUserGrpcMiddleware(jwtService)
		},
	)
	return e
}
