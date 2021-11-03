package providers

import "go.uber.org/dig"

type ProviderInterface interface {
	Provide(container *dig.Container) error
}
