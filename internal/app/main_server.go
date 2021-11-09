package app

import (
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	providers2 "go-user-microservice/internal/app/card/providers"
	"go-user-microservice/internal/app/user"
	userProviders "go-user-microservice/internal/app/user/providers"
	wallet2 "go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/app/wallet/middlewares"
	walletProvider "go-user-microservice/internal/app/wallet/providers"
	"go-user-microservice/internal/pkg/config"
	sharedContainers "go-user-microservice/internal/pkg/providers/containers"
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
	if e := providers2.ProvideCardRepositories(s.container); e != nil {
		return e
	}
	if e := providers2.ProvideCardServices(s.container); e != nil {
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
	e = s.container.Invoke(func(config *config.Config) {
		configuration = config
	})
	if e != nil {
		return e
	}
	e = s.container.Invoke(func(walletServer *wallet2.GrpcWalletServer) {
		walletGrpcServer = walletServer
	})
	if e != nil {
		return e
	}
	e = s.container.Invoke(func(userGrpcMiddleware *middlewares.WalletGrpcMiddleware) {
		userMiddleware = userGrpcMiddleware
	})
	if e != nil {
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
	e := s.container.Invoke(func(userServer *user.GrpcUserServer) {
		userGrpcServer = userServer
	})
	if e != nil {
		return nil, e
	}
	return userGrpcServer, nil
}

func (s *Server) GetWalletGrpcServer() (*wallet2.GrpcWalletServer, error) {
	walletGrpcServer := new(wallet2.GrpcWalletServer)
	e := s.container.Invoke(func(walletServer *wallet2.GrpcWalletServer) {
		walletGrpcServer = walletServer
	})
	if e != nil {
		return nil, e
	}
	return walletGrpcServer, nil
}
