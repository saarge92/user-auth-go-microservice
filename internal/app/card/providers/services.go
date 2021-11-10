package providers

import (
	"go-user-microservice/internal/app/card"
	stripeInterface "go-user-microservice/internal/pkg/domain/services/stripe"
	"go.uber.org/dig"
)

func ProvideCardServices(container *dig.Container) error {
	return container.Provide(
		func(
			cardRepository *card.RepositoryCard,
			cardStripeService stripeInterface.CardStripeServiceInterface,
		) *card.ServiceCard {
			return card.NewServiceCard(cardRepository, cardStripeService)
		},
	)
}
