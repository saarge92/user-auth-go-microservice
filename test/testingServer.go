package test

import (
	"github.com/joho/godotenv"
	"go-user-microservice/internal/app/providers/containers"
	"go-user-microservice/internal/app/server"
	"go.uber.org/dig"
	"os"
	"path"
	"runtime"
)

type ServerTest struct {
	container *dig.Container
}

func NewServerTest() *ServerTest {
	serverTest := &ServerTest{}
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
	e := containers.ProvideConnections(s.container)
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

func (s *ServerTest) Start() error {
	return nil
}

func (s *ServerTest) GetDIContainer() *dig.Container {
	return s.container
}

func (s *ServerTest) GetUserGrpcServer() (*server.UserGrpcServer, error) {
	userGrpcServer := new(server.UserGrpcServer)
	e := s.container.Invoke(func(userServer *server.UserGrpcServer) {
		userGrpcServer = userServer
	})
	if e != nil {
		return nil, e
	}
	return userGrpcServer, nil
}
