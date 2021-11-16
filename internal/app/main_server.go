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
	container        *dig.Container
	userGrpcServer   *user.GrpcUserServer
	walletGrpcServer *walletApp.GrpcWalletServer
	cardGrpcServer   *card.GrpcServerCard
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
	if e := cardProvider.ProvideCardMiddleware(s.container); e != nil {
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

func (s *Server) initUserGrpcServer() error {
	if e := s.container.Invoke(func(userServer *user.GrpcUserServer) {
		s.userGrpcServer = userServer
	}); e != nil {
		return e
	}
	return nil
}

func (s *Server) initWalletGrpcServer() error {
	if e := s.container.Invoke(func(userServer *walletApp.GrpcWalletServer) {
		s.walletGrpcServer = userServer
	}); e != nil {
		return e
	}
	return nil
}

func (s *Server) initCardGrpcServer() error {
	if e := s.container.Invoke(func(cardServer *card.GrpcServerCard) {
		s.cardGrpcServer = cardServer
	}); e != nil {
		return e
	}
	return nil
}

func (s *Server) Start() error {
	if e := s.initUserGrpcServer(); e != nil {
		return e
	}
	if e := s.initWalletGrpcServer(); e != nil {
		return e
	}
	if e := s.initCardGrpcServer(); e != nil {
		return e
	}
	return s.runAndRegisterServers()
}

func (s *Server) runAndRegisterServers() error {
	var configuration *config.Config
	var userMiddleware *middlewares.WalletGrpcMiddleware
	var cardMiddleware *card.GrpcCardMiddleware
	if e := s.container.Invoke(func(config *config.Config) {
		configuration = config
	}); e != nil {
		return e
	}
	if e := s.container.Invoke(func(userGrpcMiddlewareInstance *middlewares.WalletGrpcMiddleware) {
		userMiddleware = userGrpcMiddlewareInstance
	}); e != nil {
		return e
	}
	if e := s.container.Invoke(func(cardGrpcMiddlewareInstance *card.GrpcCardMiddleware) {
		cardMiddleware = cardGrpcMiddlewareInstance
	}); e != nil {
		return e
	}
	serv := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcLogrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			userMiddleware.CreateWalletAuthenticated,
			cardMiddleware.CreateCardAuthenticated,
		),
	)

	user_server.RegisterUserServiceServer(serv, s.userGrpcServer)
	wallet.RegisterWalletServiceServer(serv, s.walletGrpcServer)
	cardGrpc.RegisterCardServiceServer(serv, s.cardGrpcServer)
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

func (s *Server) GetUserGrpcServer() *user.GrpcUserServer {
	return s.userGrpcServer
}

func (s *Server) GetWalletGrpcServer() *walletApp.GrpcWalletServer {
	return s.walletGrpcServer
}

func (s *Server) GetCardGrpcServer() *card.GrpcServerCard {
	return s.cardGrpcServer
}
