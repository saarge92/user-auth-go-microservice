package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/user/forms/builders"
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/pkg/domain/providers"
)

type GrpcServerProvider struct {
	wallet *wallet.GrpcWalletServer
	user   *user.GrpcUserServer
	card   *card.GrpcServerCard
}

func NewGrpcServerProvider(
	serviceProvider providers.ServiceProviderInterface,
) *GrpcServerProvider {
	walletGrpcServer := wallet.NewWalletGrpcServer(
		serviceProvider.WalletService(),
	)
	userGrpcServer := user.NewUserGrpcServer(
		&builders.UserFormBuilder{},
		serviceProvider.AuthService(),
	)
	cardGrpcServer := card.NewGrpcServerCard(
		&forms.CardFormBuilder{},
		serviceProvider.CardService(),
	)
	return &GrpcServerProvider{
		wallet: walletGrpcServer,
		user:   userGrpcServer,
		card:   cardGrpcServer,
	}
}

func (p *GrpcServerProvider) UserGrpcServer() *user.GrpcUserServer {
	return p.user
}

func (p *GrpcServerProvider) WalletGrpcServer() *wallet.GrpcWalletServer {
	return p.wallet
}

func (p *GrpcServerProvider) CardGrpcServer() *card.GrpcServerCard {
	return p.card
}
