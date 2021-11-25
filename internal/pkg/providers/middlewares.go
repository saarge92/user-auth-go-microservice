package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/pkg/domain/providers"
)

type GrpcMiddlewareProvider struct {
	wallet *wallet.GrpcWalletMiddleware
	card   *card.GrpcCardMiddleware
}

func NewGrpcMiddlewareProvider(
	serviceProvider providers.ServiceProviderInterface,
) *GrpcMiddlewareProvider {
	walletGrpcMiddleware := wallet.NewWalletGrpcServerMiddleware(serviceProvider.UserAuthContextService())
	cardGrpcMiddleware := card.NewGrpcCardMiddleware(serviceProvider.UserAuthContextService())
	return &GrpcMiddlewareProvider{
		wallet: walletGrpcMiddleware,
		card:   cardGrpcMiddleware,
	}
}

func (p *GrpcMiddlewareProvider) Wallet() *wallet.GrpcWalletMiddleware {
	return p.wallet
}

func (p *GrpcMiddlewareProvider) Card() *card.GrpcCardMiddleware {
	return p.card
}
