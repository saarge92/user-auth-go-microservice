package app

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/providers/containers"
	"go-user-microservice/internal/server"
	"go-user-microservice/pkg/protobuf/user"
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
func (s *Server) InitContainer(driverDB string) error {
	e := containers.ProvideConnections(s.container, driverDB)
	if e != nil {
		return e
	}
	e = containers.ProvideRepositories(s.container)
	if e != nil {
		return e
	}
	e = containers.ProvideUserServices(s.container)
	if e != nil {
		return e
	}
	e = containers.ProvideForms(s.container)
	if e != nil {
		return e
	}
	e = containers.ProvideGrpcServers(s.container)
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
	user.RegisterUserServiceServer(serv, userGrpcServer)
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
