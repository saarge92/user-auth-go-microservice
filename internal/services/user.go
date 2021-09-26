package services

import (
	repositories "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/errorLists"
	"go-user-microservice/internal/forms"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userRepository     repositories.UserRepository
	userRemoteServices RemoteUserService
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) SignUp(form *forms.SignUp, chanResp chan<- interface{}) (*entites.User, error) {
	userExist, e := s.userRemoteServices.CheckRemoteUser(form.Login)
	if e != nil {
		return nil, e
	}
	if !userExist {
		return nil, status.Error(codes.NotFound, errorLists.UserNotFoundOnRemote)
	}
	user := &entites.User{}
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(form.Password), 14)
	if e != nil {
		chanResp <- e
		return nil, e
	}
	user.Password = string(passwordHash)
	user.Login = form.Login
	user.Name = form.Name
	if e = s.userRepository.Create(user); e != nil {
		chanResp <- nil
		return nil, e
	}
	chanResp <- nil
	return user, nil
}
