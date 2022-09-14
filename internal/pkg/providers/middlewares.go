package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/payment"
	"go-user-microservice/internal/app/wallet"
)

type GrpcMiddlewareProvider struct {
	wallet  *wallet.GrpcWalletMiddleware
	card    *card.GrpcCardMiddleware
	payment *payment.GrpcPaymentMiddleware
}

func NewGrpcMiddlewareProvider(
	serviceProvider *ServiceProvider,
) *GrpcMiddlewareProvider {
	walletGrpcMiddleware := wallet.NewWalletGrpcServerMiddleware(serviceProvider.UserAuthContextService())
	cardGrpcMiddleware := card.NewGrpcCardMiddleware(serviceProvider.UserAuthContextService())
	paymentGrpcMiddleware := payment.NewGrpcPaymentMiddleware(serviceProvider.UserAuthContextService())
	return &GrpcMiddlewareProvider{
		wallet:  walletGrpcMiddleware,
		card:    cardGrpcMiddleware,
		payment: paymentGrpcMiddleware,
	}
}

func (p *GrpcMiddlewareProvider) Wallet() *wallet.GrpcWalletMiddleware {
	return p.wallet
}

func (p *GrpcMiddlewareProvider) Card() *card.GrpcCardMiddleware {
	return p.card
}

func (p *GrpcMiddlewareProvider) Payment() *payment.GrpcPaymentMiddleware {
	return p.payment
}
