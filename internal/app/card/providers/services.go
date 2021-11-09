package providers

import (
	"go-user-microservice/internal/app/card"
	"go.uber.org/dig"
)

func ProvideCardServices(container *dig.Container) error {
	return container.Provide(
		func(cardRepository *card.RepositoryCard) *card.ServiceCard {
			return card.NewServiceCard(cardRepository)
		})
}
