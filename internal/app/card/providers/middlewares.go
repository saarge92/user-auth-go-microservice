package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/pkg/services"
	"go.uber.org/dig"
)

func ProvideCardMiddleware(container *dig.Container) error {
	return container.Provide(
		func(authUserContextService *services.UserAuthContextService) *card.GrpcCardMiddleware {
			return card.NewGrpcCardMiddleware(authUserContextService)
		},
	)
}
