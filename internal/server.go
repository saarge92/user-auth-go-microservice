package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/providers"
	"go-user-microservice/internal/providers/functions"
	"go-user-microservice/internal/repositories"
	"go-user-microservice/internal/server"
	"go-user-microservice/pkg/protobuf/user"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Container *dig.Container
}

func NewServer() *Server {
	mainServer := &Server{}
	e := mainServer.initContainer()
	if e != nil {
		log.Fatal(e)
	}
	return mainServer
}

func (s *Server) initContainer() error {
	s.Container = dig.New()
	e := s.Container.Provide(func() *config.Config {
		return config.NewConfig()
	})
	if e != nil {
		return e
	}
	e = s.Container.Provide(func(config *config.Config) *providers.ConnectionProvider {
		return providers.NewConnectionProvider(config)
	})
	if e != nil {
		return e
	}
	e = s.Container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.UserRepository {
			return repositories.NewUserRepository(connProvider.GetCoreConnection())
		})
	if e != nil {
		return e
	}
	e = functions.ProvideUserRepositories(s.Container)
	if e != nil {
		return e
	}
	e = functions.ProvideGrpcServers(s.Container)
	if e != nil {
		return e
	}
	return nil
}

func (s *Server) Start() error {
	serv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpcLogrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
		),
	)
	var userGrpcServer *server.UserGrpcServer
	var configuration *config.Config
	e := s.Container.Invoke(func(userServer *server.UserGrpcServer) {
		userGrpcServer = userServer
	})
	if e != nil {
		return e
	}
	e = s.Container.Invoke(func(config *config.Config) {
		configuration = config
	})
	if e != nil {
		return e
	}
	user.RegisterUserServiceServer(serv, userGrpcServer)
	listener, e := net.Listen("tcp", fmt.Sprintf(":%s", configuration.GrpcPort))
	if e != nil {
		return e
	}
	log.Debug("User microservice server running on port", configuration.GrpcPort)
	return serv.Serve(listener)
}
