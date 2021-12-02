package providers

import (
	cardServices "go-user-microservice/internal/app/card/services"
	"go-user-microservice/internal/app/payment/domain"
	"go-user-microservice/internal/app/user/services"
	domainService "go-user-microservice/internal/pkg/domain/services"
	"go-user-microservice/internal/pkg/domain/services/stripe"
)

type ServiceProvider interface {
	AuthService() *services.AuthService
	JwtService() *services.JwtService
	RemoteUserService() domainService.RemoteUserService
	UserService() *services.ServiceUser
	WalletService() domainService.WalletService
	StripeAccountService() stripe.AccountStripeService
	StripeCardService() stripe.CardStripeService
	UserAuthContextService() *services.UserAuthContextService
	CardService() *cardServices.ServiceCard
	PaymentService() domain.PaymentService
}
