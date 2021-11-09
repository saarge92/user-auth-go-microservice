package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/pkg/providers"
	"go.uber.org/dig"
)

func ProvideCardRepositories(container *dig.Container) error {
	return container.Provide(
		func(connectionProvider *providers.ConnectionProvider) *card.RepositoryCard {
			return card.NewRepositoryCard(connectionProvider.GetCoreConnection())
		},
	)
}
