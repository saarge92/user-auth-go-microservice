package providers

import (
	userServices "go-user-microservice/internal/app/user/services"
	walletServices "go-user-microservice/internal/app/wallet/services"
	"go-user-microservice/internal/pkg/config"
	userDomain "go-user-microservice/internal/pkg/domain/services"
	stripeDomain "go-user-microservice/internal/pkg/domain/services/stripe"
	"go-user-microservice/internal/pkg/services"
	"go-user-microservice/internal/pkg/services/stripe"
)

type ServiceProvider struct {
	accountStripeService   stripeDomain.AccountStripeServiceInterface
	cardStripeService      stripeDomain.CardStripeServiceInterface
	userRemoteService      userDomain.RemoteUserServiceInterface
	userService            *userServices.ServiceUser
	jwtAuthService         *userServices.JwtService
	authService            *userServices.AuthService
	userAuthContextService *services.UserAuthContextService
	walletService          *walletServices.WalletService
}

func NewServiceProvider(
	config *config.Config,
	repositoryProvider *RepositoryProvider,
	dbConnectionProvider *DatabaseConnectionProvider,
	stripeClientProvider *ClientStripeProvider,
) *ServiceProvider {
	// stripe
	accountStripeService := stripe.NewAccountStripeService(stripeClientProvider.MainClient())
	cardStripeService := stripe.NewCardStripeService(stripeClientProvider.MainClient())

	// user
	remoteUserService := userServices.NewRemoteUserService(config)
	userService := userServices.NewUserService(
		repositoryProvider.UserRepository(),
		repositoryProvider.CountryRepository(),
		remoteUserService,
		accountStripeService,
	)
	jwtService := userServices.NewJwtService(
		config,
		repositoryProvider.UserRepository(),
	)
	authService := userServices.NewAuthService(userService, jwtService)

	//wallet
	walletService := walletServices.NewWalletService(
		repositoryProvider.WalletRepository(),
		repositoryProvider.UserRepository(),
		repositoryProvider.CurrencyRepository(),
		dbConnectionProvider.GetCoreConnection(),
	)
	userAuthContextService := services.NewUserAuthContextService(jwtService)
	return &ServiceProvider{
		accountStripeService:   accountStripeService,
		cardStripeService:      cardStripeService,
		userRemoteService:      remoteUserService,
		userService:            userService,
		jwtAuthService:         jwtService,
		authService:            authService,
		walletService:          walletService,
		userAuthContextService: userAuthContextService,
	}
}

func (p *ServiceProvider) AuthService() *userServices.AuthService {
	return p.authService
}

func (p *ServiceProvider) JwtService() *userServices.JwtService {
	return p.jwtAuthService
}

func (p *ServiceProvider) RemoteUserService() userDomain.RemoteUserServiceInterface {
	return p.userRemoteService
}

func (p *ServiceProvider) UserService() *userServices.ServiceUser {
	return p.userService
}

func (p *ServiceProvider) WalletService() userDomain.WalletServiceInterface {
	return p.walletService
}

func (p *ServiceProvider) StripeAccountService() stripeDomain.AccountStripeServiceInterface {
	return p.accountStripeService
}

func (p *ServiceProvider) StripeCardService() stripeDomain.CardStripeServiceInterface {
	return p.cardStripeService
}

func (p *ServiceProvider) UserAuthContextService() *services.UserAuthContextService {
	return p.userAuthContextService
}
