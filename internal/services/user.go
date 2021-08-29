package services

import (
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/forms"
	repositories "go-user-microservice/internal/interfaces/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) SignUp(form *forms.SignUp) (chan *entites.User, chan error) {

	return nil, nil
}
