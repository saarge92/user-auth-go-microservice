package test

import (
	"github.com/joho/godotenv"
	containers2 "go-user-microservice/internal/app/providers/containers"
	"go.uber.org/dig"
	"os"
	"path"
	"runtime"
)

type ServerTest struct {
	container *dig.Container
}

func NewServerTest() *ServerTest {
	server := &ServerTest{}
	server.container = dig.New()
	return server
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
	return containers2.ProvideConfig(s.container)
}

func (s *ServerTest) InitContainer() error {
	e := containers2.ProvideConnections(s.container)
	if e != nil {
		return e
	}
	e = containers2.ProvideRepositories(s.container)
	if e != nil {
		return e
	}
	e = containers2.ProvideUserServices(s.container)
	if e != nil {
		return e
	}
	e = containers2.ProvideForms(s.container)
	if e != nil {
		return e
	}
	e = containers2.ProvideGrpcServers(s.container)
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
