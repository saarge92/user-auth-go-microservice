package test_providers

import (
	stripeInterface "go-user-microservice/internal/pkg/domain/services/stripe"
	serviceMock "go-user-microservice/test/services"
	"go.uber.org/dig"
)

func ProvideStripe(container *dig.Container) error {
	e := container.Provide(
		func() stripeInterface.AccountStripeServiceInterface {
			mock := &serviceMock.AccountStripeServiceMock{}
			var mockInterface stripeInterface.AccountStripeServiceInterface = mock
			return mockInterface
		},
	)
	return e
}
