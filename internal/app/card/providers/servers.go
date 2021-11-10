package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/card/forms"
	"go.uber.org/dig"
)

func ProvideCardServer(s *dig.Container) error {
	return s.Provide(
		func(serviceCard *card.ServiceCard) *card.GrpcServerCard {
			cardFormBuilder := &forms.CardFormBuilder{}
			return card.NewGrpcServerCard(cardFormBuilder, serviceCard)
		},
	)
}
