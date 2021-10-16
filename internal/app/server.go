package app

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/app/config"
	"go-user-microservice/internal/app/middlewares"
	"go-user-microservice/internal/app/providers/containers"
	"go-user-microservice/internal/app/server"
	"go-user-microservice/pkg/protobuf/user"
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
	if e := containers.ProvideConnections(s.container); e != nil {
		return e
	}
	if e := containers.ProvideRepositories(s.container); e != nil {
		return e
	}
	if e := containers.ProvideUserServices(s.container); e != nil {
		return e
	}
	if e := containers.ProvideWalletServices(s.container); e != nil {
		return e
	}
	if e := containers.ProvideUserMiddlewares(s.container); e != nil {
		return e
	}
	if e := containers.ProvideForms(s.container); e != nil {
		return e
	}
	if e := containers.ProvideGrpcServers(s.container); e != nil {
		return e
	}
	return nil
}

func (s *Server) Start() error {
	var userGrpcServer *server.UserGrpcServer
	var walletGrpcServer *server.WalletGrpcServer
	var configuration *config.Config
	var userMiddleware *middlewares.UserGrpcMiddleware
	e := s.container.Invoke(func(userServer *server.UserGrpcServer) {
		userGrpcServer = userServer
	})
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
		grpc_middleware.WithUnaryServerChain(
			grpcLogrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			userMiddleware.IsAuthenticatedMiddleware,
		),
	)
	user.RegisterUserServiceServer(serv, userGrpcServer)
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
