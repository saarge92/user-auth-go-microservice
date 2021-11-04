package test_providers

import (
	serviceMock "go-user-microservice/test/services"
	"go.uber.org/dig"
)

func ProvideStripe(container *dig.Container) error {
	e := container.Provide(
		func() *serviceMock.AccountStripeServiceMock {
			return &serviceMock.AccountStripeServiceMock{}
		},
	)
	return e
}
