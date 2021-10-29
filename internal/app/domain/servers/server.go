package servers

import (
	"go-user-microservice/internal/app/server"
	"go.uber.org/dig"
)

type ServerInterface interface {
	InitConfig() error
	InitContainer() error
	Start() error
	GetDIContainer() *dig.Container
	GetUserGrpcServer() (*server.UserGrpcServer, error)
	GetWalletGrpcServer() (*server.WalletGrpcServer, error)
}
