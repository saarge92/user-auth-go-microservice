package servers

import "go.uber.org/dig"

type ServerInterface interface {
	InitConfig() error
	InitContainer() error
	Start() error
	GetDIContainer() *dig.Container
}
