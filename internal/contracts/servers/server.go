package servers

import "go.uber.org/dig"

type ServerInterface interface {
	InitConfig() error
	InitContainer(driverDB string) error
	Start() error
	GetDIContainer() *dig.Container
}
