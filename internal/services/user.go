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
	userExist, e := s.userRemoteServices.CheckRemoteUser(form.Inn)
	if e != nil {
		chanResp <- e
		return nil, e
	}
	if !userExist {
		chanResp <- nil
		return nil, status.Error(codes.NotFound, errorlists.UserNotFoundOnRemote)
	}
	user := &entites.User{}
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if e != nil {
		chanResp <- e
		return nil, e
	}
	user.Password = string(passwordHash)
	user.Login = form.Login
	user.Name = form.Name
	user.Inn = form.Inn
	if e = s.userRepository.Create(user); e != nil {
		chanResp <- e
		return nil, e
	}
	chanResp <- nil
	return user, nil
}

func (s *UserService) SignIn(
	form *forms.SignIn,
	chanResp chan<- interface{},
) (*entites.User, error) {
	user, e := s.userRepository.GetUser(form.Login)
	if e != nil {
		chanResp <- e
		return nil, e
	}
	if user == nil {
		userErr := status.Error(codes.NotFound, errorlists.UserNotFound)
		chanResp <- userErr
		return nil, userErr
	}
	e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if e != nil {
		chanResp <- e
		return nil, e
	}
	chanResp <- nil
	return user, nil
}
