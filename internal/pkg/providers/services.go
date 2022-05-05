package providers

import (
	cardServices "go-user-microservice/internal/app/card/services"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/payment/services"
	userServices "go-user-microservice/internal/app/user/services"
	walletServices "go-user-microservice/internal/app/wallet/services"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/domain/providers"
	userDomain "go-user-microservice/internal/pkg/domain/services"
	stripeDomain "go-user-microservice/internal/pkg/domain/services/stripe"
)

type ServiceProvider struct {
	accountStripeService   stripeDomain.AccountStripeService
	cardStripeService      stripeDomain.CardStripeService
	userRemoteService      userDomain.RemoteUserService
	userService            *userServices.User
	jwtAuthService         *userServices.JwtService
	authService            *userServices.Auth
	userAuthContextService *userServices.UserAuthContextService
	walletService          *walletServices.WalletService
	cardService            *cardServices.ServiceCard
	paymentService         domain.PaymentService
}

func NewServiceProvider(
	config *config.Config,
	repositoryProvider providers.RepositoryProvider,
	dbConnectionProvider providers.DatabaseConnectionProvider,
	stripeServiceProvider providers.StripeServiceProvider,
) *ServiceProvider {
	// user
	remoteUserService := userServices.NewRemoteUserService(config)
	userService := userServices.NewUserService(
		repositoryProvider.UserRepository(),
		repositoryProvider.CountryRepository(),
		remoteUserService,
		stripeServiceProvider.Account(),
		repositoryProvider.RoleRepository(),
	)
	jwtService := userServices.NewJwtService(
		config,
		repositoryProvider.UserRepository(),
	)
	authService := userServices.NewAuthService(userService, jwtService)

	// wallet
	walletService := walletServices.NewWalletService(
		repositoryProvider.WalletRepository(),
		repositoryProvider.UserRepository(),
		repositoryProvider.CurrencyRepository(),
		dbConnectionProvider.GetCoreConnection(),
	)
	userAuthContextService := userServices.NewUserAuthContextService(jwtService)
	cardService := cardServices.NewServiceCard(repositoryProvider.CardRepository(), stripeServiceProvider.Card())
	paymentService := services.NewPaymentService(
		repositoryProvider.OperationStory(),
		repositoryProvider.WalletRepository(),
		repositoryProvider.CardRepository(),
		stripeServiceProvider.Charge(),
		dbConnectionProvider.GetCoreConnection(),
	)
	return &ServiceProvider{
		accountStripeService:   stripeServiceProvider.Account(),
		cardStripeService:      stripeServiceProvider.Card(),
		userRemoteService:      remoteUserService,
		userService:            userService,
		jwtAuthService:         jwtService,
		authService:            authService,
		walletService:          walletService,
		userAuthContextService: userAuthContextService,
		cardService:            cardService,
		paymentService:         paymentService,
	}
}

func (p *ServiceProvider) AuthService() *userServices.Auth {
	return p.authService
}

func (p *ServiceProvider) JwtService() *userServices.JwtService {
	return p.jwtAuthService
}

func (p *ServiceProvider) RemoteUserService() userDomain.RemoteUserService {
	return p.userRemoteService
}

func (p *ServiceProvider) UserService() *userServices.User {
	return p.userService
}

func (p *ServiceProvider) WalletService() userDomain.WalletService {
	return p.walletService
}

func (p *ServiceProvider) StripeAccountService() stripeDomain.AccountStripeService {
	return p.accountStripeService
}

func (p *ServiceProvider) StripeCardService() stripeDomain.CardStripeService {
	return p.cardStripeService
}

func (p *ServiceProvider) UserAuthContextService() *userServices.UserAuthContextService {
	return p.userAuthContextService
}

func (p *ServiceProvider) CardService() *cardServices.ServiceCard {
	return p.cardService
}

func (p *ServiceProvider) PaymentService() domain.PaymentService {
	return p.paymentService
}
