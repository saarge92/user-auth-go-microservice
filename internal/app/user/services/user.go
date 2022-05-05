package services

import (
	"context"
	"database/sql"
	"go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/user/forms"
	"go-user-microservice/internal/app/user/repositories"
	sharedRepoInterfaces "go-user-microservice/internal/pkg/domain/repositories"
	userDomain "go-user-microservice/internal/pkg/domain/services"
	stripeDomain "go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/entites"
	"go-user-microservice/internal/pkg/errorlists"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	userRepository       domain.UserRepository
	countryRepository    sharedRepoInterfaces.CountryRepository
	userRemoteServices   userDomain.RemoteUserService
	stripeAccountService stripeDomain.AccountStripeService
	userRolesRepository  *repositories.Role
}

func NewUserService(
	userRepository domain.UserRepository,
	countryRepository sharedRepoInterfaces.CountryRepository,
	userRemoteService userDomain.RemoteUserService,
	stripeAccountService stripeDomain.AccountStripeService,
	userRolesRepository *repositories.Role,
) *User {
	return &User{
		userRepository:       userRepository,
		userRemoteServices:   userRemoteService,
		countryRepository:    countryRepository,
		stripeAccountService: stripeAccountService,
		userRolesRepository:  userRolesRepository,
	}
}

func (s *User) SignUp(form *forms.SignUp) (*entities.User, error) {
	var country *entites.Country
	var checkError error
	if country, checkError = s.checkUserDataWithCountryResponse(form); checkError != nil {
		return nil, checkError
	}

	stripeAccountData := &dto.StripeAccountCreate{
		Email:        form.Login,
		Country:      form.Country,
		CardPayments: true,
		Transfers:    true,
	}
	userEntity := &entities.User{}

	if country != nil {
		userEntity.CountryID = sql.NullInt64{Int64: int64(country.ID), Valid: true}
		stripeAccountData.Country = country.CodeTwo
	}
	accountResponse, customerResponse, e := s.stripeAccountService.Create(stripeAccountData)
	if e != nil {
		return nil, e
	}
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if e != nil {
		return nil, e
	}
	userEntity.Password = string(passwordHash)
	userEntity.Login = form.Login
	userEntity.Name = form.Name
	userEntity.Inn = form.Inn
	userEntity.AccountProviderID = accountResponse.ID
	userEntity.CustomerProviderID = customerResponse.ID
	if e = s.userRepository.Create(userEntity); e != nil {
		return nil, e
	}
	if e = s.userRolesRepository.AddUserToRole(context.Background(), userEntity.ID, entities.UserRoleID); e != nil {
		return nil, e
	}
	return userEntity, nil
}

func (s *User) SignIn(form *forms.SignIn) (*entities.User, error) {
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

func (s *User) checkUserDataWithCountryResponse(form *forms.SignUp) (*entites.Country, error) {
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
	if form.Country != "" {
		var countryError error
		var country *entites.Country
		country, countryError = s.countryRepository.GetByCodeTwo(context.Background(), form.Country)
		if countryError != nil {
			return nil, countryError
		}
		return country, nil
	}
	return nil, nil
}
