package test

import (
	"github.com/joho/godotenv"
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/user"
	userProviders "go-user-microservice/internal/app/user/providers"
	"go-user-microservice/internal/app/wallet"
	walletProviders "go-user-microservice/internal/app/wallet/providers"
	"go-user-microservice/internal/pkg/domain/providers"
	sharedContainers "go-user-microservice/internal/pkg/providers/containers"
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
	return sharedContainers.ProvideConfig(s.container)
}

func (s *ServerTest) InitContainer() error {
	if e := sharedContainers.ProvideEncryptionService(s.container); e != nil {
		return e
	}
	if e := sharedContainers.ProvideConnection(s.container); e != nil {
		return e
	}
	if e := sharedContainers.ProvideShareRepositories(s.container); e != nil {
		return e
	}
	if e := userProviders.ProviderUserRepository(s.container); e != nil {
		return e
	}
	if e := walletProviders.ProvideWalletRepository(s.container); e != nil {
		return e
	}
	if s.stripeProvideFunction != nil {
		if e := s.stripeProvideFunction(s.container); e != nil {
			return e
		}
	} else {
		if e := sharedContainers.ProvideStripeService(s.container); e != nil {
			return e
		}
	}
	if e := userProviders.ProvideUserServices(s.container); e != nil {
		return e
	}
	if e := walletProviders.ProvideWalletServices(s.container); e != nil {
		return e
	}
	if e := userProviders.ProvideGrpcMiddleware(s.container); e != nil {
		return e
	}
	if e := userProviders.ProvideUserGrpcServers(s.container); e != nil {
		return e
	}
	if e := walletProviders.ProvideWalletGrpcServer(s.container); e != nil {
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

func (s *ServerTest) GetUserGrpcServer() (*user.GrpcUserServer, error) {
	var userGrpcServer *user.GrpcUserServer
	e := s.container.Invoke(func(userServer *user.GrpcUserServer) {
		userGrpcServer = userServer
	})
	if e != nil {
		return nil, e
	}
	return userGrpcServer, nil
}

func (s *ServerTest) GetWalletGrpcServer() (*wallet.GrpcWalletServer, error) {
	var walletGrpcServer *wallet.GrpcWalletServer
	e := s.container.Invoke(func(walletServer *wallet.GrpcWalletServer) {
		walletGrpcServer = walletServer
	})
	if e != nil {
		return nil, e
	}
	return walletGrpcServer, nil
}

func (s *ServerTest) GetCardGrpcServer() (*card.GrpcServerCard, error) {
	cardGrpcServer := new(card.GrpcServerCard)
	if e := s.container.Invoke(func(cardServer *card.GrpcServerCard) {
		cardGrpcServer = cardServer
	}); e != nil {
		return nil, e
	}
	return cardGrpcServer, nil
}
