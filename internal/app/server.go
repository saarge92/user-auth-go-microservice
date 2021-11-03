package app

import (
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/middlewares"
	"go-user-microservice/internal/app/providers/containers"
	"go-user-microservice/internal/app/server"
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
	return containers.ProvideConfig(s.container)
}
func (s *Server) InitContainer() error {
	userServiceProvider := &containers.UserServicesProvider{}
	repositoryProvider := &containers.RepositoryProvider{}
	encryptionServiceProvider := &containers.EncryptionProvider{}
	connectionProvider := &containers.ConnectionProvider{}
	walletServiceProvider := &containers.WalletServiceProvider{}
	userGRPCMiddlewareProvider := &containers.UserGRPCMiddlewareProvider{}
	grpcServerProvider := &containers.GrpcServerProvider{}
	stripeServiceProvider := &containers.StripeServiceProvider{}

	if e := encryptionServiceProvider.Provide(s.container); e != nil {
		return e
	}
	if e := connectionProvider.Provide(s.container); e != nil {
		return e
	}
	if e := repositoryProvider.Provide(s.container); e != nil {
		return e
	}
	if e := stripeServiceProvider.Provide(s.container); e != nil {
		return e
	}
	if e := userServiceProvider.Provide(s.container); e != nil {
		return e
	}
	if e := walletServiceProvider.Provide(s.container); e != nil {
		return e
	}
	if e := userGRPCMiddlewareProvider.Provide(s.container); e != nil {
		return e
	}
	if e := grpcServerProvider.Provide(s.container); e != nil {
		return e
	}
	return nil
}

func (s *Server) Start() error {
	var configuration *config.Config
	var userMiddleware *middlewares.UserGrpcMiddleware
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
	e = s.container.Invoke(func(walletServer *server.WalletGrpcServer) {
		walletGrpcServer = walletServer
	})
	if e != nil {
		return e
	}
	e = s.container.Invoke(func(userGrpcMiddleware *middlewares.UserGrpcMiddleware) {
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

func (s *Server) GetUserGrpcServer() (*server.UserGrpcServer, error) {
	userGrpcServer := new(server.UserGrpcServer)
	e := s.container.Invoke(func(userServer *server.UserGrpcServer) {
		userGrpcServer = userServer
	})
	if e != nil {
		return nil, e
	}
	return userGrpcServer, nil
}

func (s *Server) GetWalletGrpcServer() (*server.WalletGrpcServer, error) {
	walletGrpcServer := new(server.WalletGrpcServer)
	e := s.container.Invoke(func(walletServer *server.WalletGrpcServer) {
		walletGrpcServer = walletServer
	})
	if e != nil {
		return nil, e
	}
	return walletGrpcServer, nil
}
