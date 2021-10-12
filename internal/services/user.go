package services

import (
	repositories "go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/errorlists"
	"go-user-microservice/internal/forms"
	"go-user-microservice/internal/services/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userRepository     repositories.UserRepositoryInterface
	userRemoteServices *user.RemoteUserService
}

func NewUserService(
	userRepository repositories.UserRepositoryInterface,
	userRemoteService *user.RemoteUserService,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRemoteServices: userRemoteService,
	}
}

func (s *UserService) SignUp(form *forms.SignUp) (*entites.User, error) {
	userExist, e := s.userRemoteServices.CheckRemoteUser(form.Inn)
	if e != nil {
		return nil, e
	}
	if !userExist {
		return nil, status.Error(codes.NotFound, errorlists.UserNotFoundOnRemote)
	}
	user := &entites.User{}
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if e != nil {
		return nil, e
	}
	user.Password = string(passwordHash)
	user.Login = form.Login
	user.Name = form.Name
	user.Inn = form.Inn
	if e = s.userRepository.Create(user); e != nil {
		return nil, e
	}
	return user, nil
}

func (s *UserService) SignIn(form *forms.SignIn) (*entites.User, error) {
	user, e := s.userRepository.GetUser(form.Login)
	if e != nil {
		return nil, e
	}
	if user == nil {
		userErr := status.Error(codes.NotFound, errorlists.UserNotFound)
		return nil, userErr
	}
	e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if e != nil {
		return nil, e
	}
	return user, nil
}
