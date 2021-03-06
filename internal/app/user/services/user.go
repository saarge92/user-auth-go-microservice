package services

import (
	"context"
	"database/sql"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/user/domain"
	userDto "go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	userErrors "go-user-microservice/internal/app/user/errors"
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

func (s *User) SignUp(ctx context.Context, form *forms.SignUp) (*entities.User, error) {
	if checkUserError := s.checkUserDataExistence(ctx, form.Login, form.Inn); checkUserError != nil {
		return nil, checkUserError
	}

	stripeAccountData := &dto.StripeAccountCreate{
		Email:        form.Login,
		Country:      form.Country,
		CardPayments: true,
		Transfers:    true,
	}
	userEntity := &entities.User{}

	country, countryError := s.checkCountry(ctx, form.Country)
	if countryError != nil {
		return nil, countryError
	}
	if country != nil {
		userEntity.CountryID = sql.NullInt64{Int64: int64(country.ID), Valid: true}
		stripeAccountData.Country = country.CodeTwo
	}

	accountResponse, customerResponse, e := s.stripeAccountService.Create(stripeAccountData)
	if e != nil {
		return nil, e
	}

	if userError := s.createUser(ctx, userEntity, form, accountResponse, customerResponse); userError != nil {
		return nil, userError
	}

	return userEntity, nil
}

func (s *User) SignIn(ctx context.Context, form *forms.SignIn) (*userDto.UserRole, error) {
	userEntity, e := s.userRepository.GetUserWithRoles(ctx, form.Login)
	unAuthError := status.Error(codes.Unauthenticated, errorlists.SignInFail)
	if e != nil {
		return nil, e
	}
	if userEntity == nil {
		return nil, unAuthError
	}
	hashPasswordBytes := []byte(userEntity.User.Password)
	sourcePasswordBytes := []byte(form.Password)
	if e = bcrypt.CompareHashAndPassword(hashPasswordBytes, sourcePasswordBytes); e != nil {
		return nil, unAuthError
	}

	return userEntity, nil
}

func (s *User) checkUserDataExistence(ctx context.Context, login string, inn uint64) error {
	userExist, e := s.userRepository.UserByInnOrLoginExist(ctx, login, inn)
	if e != nil {
		return e
	}
	if userExist {
		return userErrors.UserAlreadyExistErr
	}
	userRemoteExist, e := s.userRemoteServices.CheckRemoteUser(inn)
	if e != nil {
		return e
	}
	if !userRemoteExist {
		return userErrors.RemoteInnNotFoundErr
	}

	return nil
}

func (s *User) checkCountry(ctx context.Context, countryCode string) (*entites.Country, error) {
	if countryCode != "" {
		var countryError error
		var country *entites.Country
		country, countryError = s.countryRepository.GetByCodeTwo(ctx, countryCode)
		if countryError != nil {
			return nil, countryError
		}
		return country, nil
	}
	return nil, nil
}

func (s *User) createUser(
	ctx context.Context,
	userEntity *entities.User,
	form *forms.SignUp,
	accountResponse *stripe.Account,
	customerResponse *stripe.Customer,
) error {
	passwordHash, e := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if e != nil {
		return e
	}
	userEntity.Password = string(passwordHash)
	userEntity.Login = form.Login
	userEntity.Name = form.Name
	userEntity.Inn = form.Inn
	userEntity.AccountProviderID = accountResponse.ID
	userEntity.CustomerProviderID = customerResponse.ID
	if e = s.userRepository.Create(ctx, userEntity); e != nil {
		return e
	}

	return s.userRolesRepository.AddUserToRole(context.Background(), userEntity.ID, entities.UserRoleID)
}
