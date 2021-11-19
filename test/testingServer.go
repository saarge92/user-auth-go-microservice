package test

import (
	"github.com/joho/godotenv"
	"go-user-microservice/internal/app"
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/wallet"
	sharedContainers "go-user-microservice/internal/pkg/providers/containers"
	"go-user-microservice/test/services"
	"go.uber.org/dig"
	"os"
	"path"
	"runtime"
)

type ServerTest struct {
	stripeAccountMock *services.AccountStripeServiceMock
	mainServer        *app.Server
}

func NewServerTest(
	stripeAccountMock *services.AccountStripeServiceMock,
) *ServerTest {
	serverTest := &ServerTest{
		stripeAccountMock: stripeAccountMock,
		mainServer:        &app.Server{},
	}
	return serverTest
}

func (s *ServerTest) InitConfig() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	e := os.Chdir(dir)
	if e != nil {
		panic(e)
	}
	if e := godotenv.Load(".env.test"); e != nil {
		panic(e)
	}
	return sharedContainers.ProvideConfig(s.mainServer.GetDIContainer())
}

func (s *ServerTest) InitGRPCServers() error {
	return s.mainServer.InitGRPCServers()
}

func (s *ServerTest) InitContainer() error {
	return nil
}

func (s *ServerTest) Start() error {
	return nil
}

func (s *ServerTest) GetDIContainer() *dig.Container {
	return s.mainServer.GetDIContainer()
}

func (s *ServerTest) GetUserGrpcServer() *user.GrpcUserServer {
	return s.mainServer.GetUserGrpcServer()
}

func (s *ServerTest) GetWalletGrpcServer() *wallet.GrpcWalletServer {
	return s.mainServer.GetWalletGrpcServer()
}

func (s *ServerTest) GetCardGrpcServer() *card.GrpcServerCard {
	return s.mainServer.GetCardGrpcServer()
}
