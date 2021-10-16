package member

import (
	repositories2 "go-user-microservice/internal/app/domain/repositories"
	entites2 "go-user-microservice/internal/app/entites"
	errorlists2 "go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/forms/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userRepository     repositories2.UserRepositoryInterface
	userRemoteServices *RemoteUserService
}

func NewUserService(
	userRepository repositories2.UserRepositoryInterface,
	userRemoteService *RemoteUserService,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		userRemoteServices: userRemoteService,
	}
}

func (s *UserService) SignUp(form *user.SignUp) (*entites2.User, error) {
	userExist, e := s.userRemoteServices.CheckRemoteUser(form.Inn)
	if e != nil {
		return nil, e
	}
	if !userExist {
		return nil, status.Error(codes.NotFound, errorlists2.UserNotFoundOnRemote)
	}
	userEntity := &entites2.User{}
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

func (s *UserService) SignIn(form *user.SignIn) (*entites2.User, error) {
	user, e := s.userRepository.GetUser(form.Login)
	if e != nil {
		return nil, e
	}
	if user == nil {
		userErr := status.Error(codes.NotFound, errorlists2.UserNotFound)
		return nil, userErr
	}
	e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if e != nil {
		return nil, e
	}
	return user, nil
}
