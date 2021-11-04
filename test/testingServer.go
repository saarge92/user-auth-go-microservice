package test

import (
	"github.com/joho/godotenv"
	"go-user-microservice/internal/app/domain/providers"
	"go-user-microservice/internal/app/providers/containers"
	"go-user-microservice/internal/app/server"
	"go.uber.org/dig"
	"os"
	"path"
	"runtime"
)

type ServerTest struct {
	container             *dig.Container
	stripeProvideFunction providers.ProvideFunction
}

func NewServerTest(
	stripeServiceFunction providers.ProvideFunction,
) *ServerTest {
	serverTest := &ServerTest{
		stripeProvideFunction: stripeServiceFunction,
	}
	serverTest.container = dig.New()
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
	return containers.ProvideConfig(s.container)
}

func (s *ServerTest) InitContainer() error {
	if e := containers.ProvideEncryptionService(s.container); e != nil {
		return e
	}
	if e := containers.ProvideConnection(s.container); e != nil {
		return e
	}
	if e := containers.ProvideRepositoryProvider(s.container); e != nil {
		return e
	}
	if s.stripeProvideFunction != nil {
		if e := s.stripeProvideFunction(s.container); e != nil {
			return e
		}
	} else {
		if e := containers.ProvideStripeService(s.container); e != nil {
			return e
		}
	}
	if e := containers.ProvideUserServices(s.container); e != nil {
		return e
	}
	if e := containers.ProvideWalletServices(s.container); e != nil {
		return e
	}
	if e := containers.ProvideGrpcMiddleware(s.container); e != nil {
		return e
	}
	if e := containers.ProvideGrpcServers(s.container); e != nil {
		return e
	}
	return nil
}

func (s *ServerTest) Start() error {
	return nil
}

func (s *ServerTest) GetDIContainer() *dig.Container {
	return s.container
}

func (s *ServerTest) GetUserGrpcServer() (*server.UserGrpcServer, error) {
	var userGrpcServer *server.UserGrpcServer
	e := s.container.Invoke(func(userServer *server.UserGrpcServer) {
		userGrpcServer = userServer
	})
	if e != nil {
		return nil, e
	}
	return userGrpcServer, nil
}

func (s *ServerTest) GetWalletGrpcServer() (*server.WalletGrpcServer, error) {
	var walletGrpcServer *server.WalletGrpcServer
	e := s.container.Invoke(func(walletServer *server.WalletGrpcServer) {
		walletGrpcServer = walletServer
	})
	if e != nil {
		return nil, e
	}
	return walletGrpcServer, nil
}
