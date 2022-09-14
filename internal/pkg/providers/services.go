package providers

import (
	cardServices "go-user-microservice/internal/app/card/services"
	paymentServices "go-user-microservice/internal/app/payment/services"
	userServices "go-user-microservice/internal/app/user/services"
	walletServices "go-user-microservice/internal/app/wallet/services"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/services/stripe"
)

type ServiceProvider struct {
	accountStripeService   *stripe.AccountStripeService
	cardStripeService      *stripe.CardStripeService
	userRemoteService      *userServices.RemoteUser
	userService            *userServices.User
	jwtAuthService         *userServices.JwtService
	authService            *userServices.Auth
	userAuthContextService *userServices.UserAuthContextService
	walletService          *walletServices.WalletService
	cardService            *cardServices.ServiceCard
	paymentService         *paymentServices.Payment
}

func NewServiceProvider(
	config *config.Config,
	repositoryProvider *RepositoryProvider,
	stripeServiceProvider *StripeServiceProvider,
) *ServiceProvider {
	userRemoteService := userServices.NewRemoteUserService(config)
	userService := userServices.NewUserService(
		repositoryProvider.UserRepository(),
		repositoryProvider.CountryRepository(),
		userRemoteService,
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
		repositoryProvider.CurrencyRepository(),
	)
	userAuthContextService := userServices.NewUserAuthContextService(jwtService)
	cardService := cardServices.NewServiceCard(repositoryProvider.CardRepository(), stripeServiceProvider.Card())
	paymentService := paymentServices.NewPaymentService(
		repositoryProvider.OperationStory(),
		repositoryProvider.WalletRepository(),
		repositoryProvider.CardRepository(),
		stripeServiceProvider.Charge(),
	)
	return &ServiceProvider{
		accountStripeService:   stripeServiceProvider.Account(),
		cardStripeService:      stripeServiceProvider.Card(),
		userRemoteService:      userRemoteService,
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

func (p *ServiceProvider) RemoteUserService() *userServices.RemoteUser {
	return p.userRemoteService
}

func (p *ServiceProvider) UserService() *userServices.User {
	return p.userService
}

func (p *ServiceProvider) WalletService() *walletServices.WalletService {
	return p.walletService
}

func (p *ServiceProvider) StripeAccountService() *stripe.AccountStripeService {
	return p.accountStripeService
}

func (p *ServiceProvider) StripeCardService() *stripe.CardStripeService {
	return p.cardStripeService
}

func (p *ServiceProvider) UserAuthContextService() *userServices.UserAuthContextService {
	return p.userAuthContextService
}

func (p *ServiceProvider) CardService() *cardServices.ServiceCard {
	return p.cardService
}

func (p *ServiceProvider) PaymentService() *paymentServices.Payment {
	return p.paymentService
}
