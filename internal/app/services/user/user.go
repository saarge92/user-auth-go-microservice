package user

import (
	"go-user-microservice/internal/app/domain/repositories"
	"go-user-microservice/internal/app/entites"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/forms/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceUser struct {
	userRepository     repositories.UserRepositoryInterface
	userRemoteServices *RemoteUserService
}

func NewUserService(
	userRepository repositories.UserRepositoryInterface,
	userRemoteService *RemoteUserService,
) *ServiceUser {
	return &ServiceUser{
		userRepository:     userRepository,
		userRemoteServices: userRemoteService,
	}
}

func (s *ServiceUser) CheckUserData(form *user.SignUp) error {
	userExist, e := s.userRepository.UserExist(form.Login)
	if e != nil {
		return e
	}
	if userExist {
		return status.Error(codes.AlreadyExists, errorlists.UserEmailAlreadyExist)
	}
	userInnExist, e := s.userRepository.UserByInnExist(form.Inn)
	if e != nil {
		return e
	}
	if userInnExist {
		return status.Error(codes.AlreadyExists, errorlists.UserInnAlreadyExist)
	}
	userRemoteExist, e := s.userRemoteServices.CheckRemoteUser(form.Inn)
	if e != nil {
		return e
	}
	if !userRemoteExist {
		return status.Error(codes.NotFound, errorlists.UserNotFoundOnRemote)
	}
	return nil
}

func (s *ServiceUser) SignUp(form *user.SignUp) (*entites.User, error) {
	if checkError := s.CheckUserData(form); checkError != nil {
		return nil, checkError
	}
	userEntity := &entites.User{}
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if e != nil {
		return nil, e
	}
	userEntity.Password = string(passwordHash)
	userEntity.Login = form.Login
	userEntity.Name = form.Name
	userEntity.Inn = form.Inn
	if e = s.userRepository.Create(userEntity); e != nil {
		return nil, e
	}
	return userEntity, nil
}

func (s *ServiceUser) SignIn(form *user.SignIn) (*entites.User, error) {
	userEntity, e := s.userRepository.GetUser(form.Login)
	unAuthError := status.Error(codes.Unauthenticated, errorlists.SignInFail)
	if e != nil {
		return nil, e
	}
	if userEntity == nil {
		return nil, unAuthError
	}
	hashPasswordBytes := []byte(userEntity.Password)
	sourcePasswordBytes := []byte(form.Password)
	if e = bcrypt.CompareHashAndPassword(hashPasswordBytes, sourcePasswordBytes); e != nil {
		return nil, unAuthError
	}

	return userEntity, nil
}
