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

type UserService struct {
	userRepository     repositories.UserRepositoryInterface
	userRemoteServices *RemoteUserService
}

func NewUserService(
	userRepository repositories.UserRepositoryInterface,
	userRemoteService *RemoteUserService,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRemoteServices: userRemoteService,
	}
}

func (s *UserService) SignUp(form *user.SignUp) (*entites.User, error) {
	userExist, e := s.userRepository.UserExist(form.Login)
	if e != nil {
		return nil, e
	}
	if userExist {
		return nil, status.Error(codes.AlreadyExists, errorlists.UserEmailAlreadyExist)
	}
	userInnExist, e := s.userRepository.UserByInnExist(form.Inn)
	if e != nil {
		return nil, e
	}
	if userInnExist {
		return nil, status.Error(codes.AlreadyExists, errorlists.UserInnAlreadyExist)
	}
	userRemoteExist, e := s.userRemoteServices.CheckRemoteUser(form.Inn)
	if e != nil {
		return nil, e
	}
	if !userRemoteExist {
		return nil, status.Error(codes.NotFound, errorlists.UserNotFoundOnRemote)
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

func (s *UserService) SignIn(form *user.SignIn) (*entites.User, error) {
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
