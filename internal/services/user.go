package services

import (
	repositories "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/errorlists"
	"go-user-microservice/internal/forms"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userRepository     repositories.UserRepository
	userRemoteServices *RemoteUserService
}

func NewUserService(
	userRepository repositories.UserRepository,
	userRemoteService *RemoteUserService,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRemoteServices: userRemoteService,
	}
}

func (s *UserService) SignUp(form *forms.SignUp, chanResp chan<- interface{}) (*entites.User, error) {
	userExist, e := s.userRemoteServices.CheckRemoteUser(1)
	if e != nil {
		chanResp <- e
		return nil, e
	}
	if !userExist {
		chanResp <- nil
		return nil, status.Error(codes.NotFound, errorlists.UserNotFoundOnRemote)
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
		chanResp <- e
		return nil, e
	}
	chanResp <- nil
	return user, nil
}
