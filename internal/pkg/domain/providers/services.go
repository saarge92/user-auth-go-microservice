package providers

import (
	services2 "go-user-microservice/internal/app/card/services"
	"go-user-microservice/internal/app/user/services"
	domainService "go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/internal/pkg/domain/services/stripe"
	sharedServices "go-user-microservice/internal/pkg/services"
)

type ServiceProviderInterface interface {
	AuthService() *services.AuthService
	JwtService() *services.JwtService
	RemoteUserService() domainService.RemoteUserServiceInterface
	UserService() *services.ServiceUser
	WalletService() domainService.WalletServiceInterface
	StripeAccountService() stripe.AccountStripeServiceInterface
	StripeCardService() stripe.CardStripeServiceInterface
	UserAuthContextService() *sharedServices.UserAuthContextService
	CardService() *services2.ServiceCard
}
