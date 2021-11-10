package app

import (
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/card"
	cardProvider "go-user-microservice/internal/app/card/providers"
	"go-user-microservice/internal/app/user"
	userProviders "go-user-microservice/internal/app/user/providers"
	walletApp "go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/app/wallet/middlewares"
	walletProvider "go-user-microservice/internal/app/wallet/providers"
	"go-user-microservice/internal/pkg/config"
	sharedContainers "go-user-microservice/internal/pkg/providers/containers"
	cardGrpc "go-user-microservice/pkg/protobuf/card"
	"go-user-microservice/pkg/protobuf/user_server"
	"go-user-microservice/pkg/protobuf/wallet"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
	"runtime"
)

type Server struct {
	container *dig.Container
}

func NewServer() *Server {
	mainServer := &Server{}
	mainServer.container = dig.New()
	return mainServer
}
func (s *Server) InitConfig() error {
	var _, filename, _, _ = runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	e := os.Chdir(dir)
	if e != nil {
		panic(e)
	}
	if e := godotenv.Load(".env"); e != nil {
		panic(e)
	}
	return sharedContainers.ProvideConfig(s.container)
}
func (s *Server) InitContainer() error {
	if e := sharedContainers.ProvideEncryptionService(s.container); e != nil {
		return e
	}
	if e := sharedContainers.ProvideConnection(s.container); e != nil {
		return e
	}
	if e := sharedContainers.ProvideShareRepositories(s.container); e != nil {
		return e
	}
	if e := sharedContainers.ProvideStripeService(s.container); e != nil {
		return e
	}
	if e := userProviders.ProviderUserRepository(s.container); e != nil {
		return e
	}
	if e := walletProvider.ProvideWalletRepository(s.container); e != nil {
		return e
	}
	if e := userProviders.ProvideUserServices(s.container); e != nil {
		return e
	}
	if e := cardProvider.ProvideCardRepositories(s.container); e != nil {
		return e
	}
	if e := cardProvider.ProvideCardServices(s.container); e != nil {
		return e
	}
	if e := walletProvider.ProvideWalletServices(s.container); e != nil {
		return e
	}
	if e := userProviders.ProvideGrpcMiddleware(s.container); e != nil {
		return e
	}
	if e := userProviders.ProvideUserGrpcServers(s.container); e != nil {
		return e
	}
	if e := walletProvider.ProvideWalletGrpcServer(s.container); e != nil {
		return e
	}
	if e := cardProvider.ProvideCardServer(s.container); e != nil {
		return e
	}
	return nil
}

func (s *Server) Start() error {
	var configuration *config.Config
	var userMiddleware *middlewares.WalletGrpcMiddleware
	userGrpcServer, e := s.GetUserGrpcServer()
	if e != nil {
		return e
	}
	walletGrpcServer, e := s.GetWalletGrpcServer()
	if e != nil {
		return e
	}
	cardGrpcServer, e := s.GetCardGrpcServer()
	if e != nil {
		return e
	}
	if e := s.container.Invoke(func(config *config.Config) {
		configuration = config
	}); e != nil {
		return e
	}
	if e = s.container.Invoke(func(walletServer *walletApp.GrpcWalletServer) {
		walletGrpcServer = walletServer
	}); e != nil {
		return e
	}
	if e := s.container.Invoke(func(userGrpcMiddleware *middlewares.WalletGrpcMiddleware) {
		userMiddleware = userGrpcMiddleware
	}); e != nil {
		return e
	}
	serv := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcLogrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			userMiddleware.IsAuthenticatedMiddleware,
		),
	)
	user_server.RegisterUserServiceServer(serv, userGrpcServer)
	wallet.RegisterWalletServiceServer(serv, walletGrpcServer)
	cardGrpc.RegisterCardServiceServer(serv, cardGrpcServer)
	listener, e := net.Listen("tcp", fmt.Sprintf(":%s", configuration.GrpcPort))
	if e != nil {
		return e
	}
	log.Debug("User microservice server running on port", configuration.GrpcPort)
	return serv.Serve(listener)
}

func (s *Server) GetDIContainer() *dig.Container {
	return s.container
}

func (s *Server) GetUserGrpcServer() (*user.GrpcUserServer, error) {
	userGrpcServer := new(user.GrpcUserServer)
	if e := s.container.Invoke(func(userServer *user.GrpcUserServer) {
		userGrpcServer = userServer
	}); e != nil {
		return nil, e
	}
	return userGrpcServer, nil
}

func (s *Server) GetWalletGrpcServer() (*walletApp.GrpcWalletServer, error) {
	walletGrpcServer := new(walletApp.GrpcWalletServer)
	if e := s.container.Invoke(func(walletServer *walletApp.GrpcWalletServer) {
		walletGrpcServer = walletServer
	}); e != nil {
		return nil, e
	}
	return walletGrpcServer, nil
}

func (s *Server) GetCardGrpcServer() (*card.GrpcServerCard, error) {
	cardGrpcServer := new(card.GrpcServerCard)
	if e := s.container.Invoke(func(cardServer *card.GrpcServerCard) {
		cardGrpcServer = cardServer
	}); e != nil {
		return nil, e
	}
	return cardGrpcServer, nil
}
