package main

import (
	"go-user-microservice/internal/config"
	repositoriesInterface "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/providers"
	"go-user-microservice/internal/repositories"
	"go-user-microservice/internal/services"
	"go.uber.org/dig"
	"log"
)

type Server struct {
	Container *dig.Container
}

func NewServer() *Server {
	server := &Server{}
	e := server.initContainer()
	if e != nil {
		log.Fatal(e)
	}
	return server
}

func (s *Server) initContainer() error {
	s.Container = dig.New()
	e := s.Container.Provide(func() *config.Config {
		return config.NewConfig()
	})
	if e != nil {
		return nil
	}
	e = s.Container.Provide(func(config *config.Config) *providers.ConnectionProvider {
		return providers.NewConnectionProvider(config)
	})
	e = s.Container.Provide(
		func(connProvider *providers.ConnectionProvider) *repositories.UserRepository {
			return repositories.NewUserRepository(connProvider.GetCoreConnection())
		})
	e = s.Container.Provide(func(userRepo *repositories.UserRepository) *services.UserService {
		var userRepositoryInterface repositoriesInterface.UserRepository = userRepo
		return services.NewUserService(userRepositoryInterface)
	})
	if e != nil {
		return e
	}
	return nil
}

func (s *Server) Start() error {
	return nil
}
