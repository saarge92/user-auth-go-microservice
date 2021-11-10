package servers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/wallet"
	"go.uber.org/dig"
)

type ServerInterface interface {
	InitConfig() error
	InitContainer() error
	Start() error
	GetDIContainer() *dig.Container
	GetUserGrpcServer() (*user.GrpcUserServer, error)
	GetWalletGrpcServer() (*wallet.GrpcWalletServer, error)
	GetCardGrpcServer() (*card.GrpcServerCard, error)
}
