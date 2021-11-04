package providers

import "go.uber.org/dig"

type ProvideFunction func(container *dig.Container) error
