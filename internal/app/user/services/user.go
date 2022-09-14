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
	"go-user-microservice/internal/app/user/request"
	"go-user-microservice/internal/pkg/dto"
	"go-user-microservice/internal/pkg/entites"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepository       domain.UserRepository
	countryRepository    domain.CountryRepository
	userRemoteServices   domain.UserRemoteService
	stripeAccountService domain.StripeAccountService
	userRolesRepository  domain.RoleRepository
}

func NewUserService(
	userRepository domain.UserRepository,
	countryRepository domain.CountryRepository,
	userRemoteService domain.UserRemoteService,
	stripeAccountService domain.StripeAccountService,
	userRolesRepository domain.RoleRepository,
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
	if e != nil {
		return nil, e
	}
	if userEntity == nil {
		return nil, userErrors.ErrSignInFail
	}
	hashPasswordBytes := []byte(userEntity.User.Password)
	sourcePasswordBytes := []byte(form.Password)
	if e = bcrypt.CompareHashAndPassword(hashPasswordBytes, sourcePasswordBytes); e != nil {
		return nil, userErrors.ErrSignInFail
	}

	return userEntity, nil
}

func (s *User) checkUserDataExistence(ctx context.Context, login string, inn uint64) error {
	userExist, e := s.userRepository.UserByInnOrLoginExist(ctx, login, inn)
	if e != nil {
		return e
	}
	if userExist {
		return userErrors.ErrUserAlreadyExists
	}
	innRequest := request.InnRequest{Inn: inn}
	userRemoteExist, e := s.userRemoteServices.CheckRemoteUser(ctx, innRequest)
	if e != nil {
		return e
	}
	if !userRemoteExist {
		return userErrors.ErrRemoteInnNotFound
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
