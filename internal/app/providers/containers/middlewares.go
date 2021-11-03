package containers

import (
	"go-user-microservice/internal/app/middlewares"
	"go-user-microservice/internal/app/services/user"
	"go.uber.org/dig"
)

type UserGRPCMiddlewareProvider struct{}

func (p *UserGRPCMiddlewareProvider) Provide(container *dig.Container) error {
	e := container.Provide(
		func(jwtService *user.JwtService) *middlewares.UserGrpcMiddleware {
			return middlewares.NewUserGrpcMiddleware(jwtService)
		},
	)
	return e
}
