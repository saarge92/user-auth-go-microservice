package providers

import (
	cardServices "go-user-microservice/internal/app/card/services"
	"go-user-microservice/internal/app/user/services"
	domainService "go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/internal/pkg/domain/services/stripe"
)

type ServiceProviderInterface interface {
	AuthService() *services.AuthService
	JwtService() *services.JwtService
	RemoteUserService() domainService.RemoteUserServiceInterface
	UserService() *services.ServiceUser
	WalletService() domainService.WalletServiceInterface
	StripeAccountService() stripe.AccountStripeServiceInterface
	StripeCardService() stripe.CardStripeServiceInterface
	UserAuthContextService() *services.UserAuthContextService
	CardService() *cardServices.ServiceCard
}
