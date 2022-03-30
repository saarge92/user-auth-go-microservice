package providers

import (
	"go-user-microservice/internal/app/card"
	"go-user-microservice/internal/app/card/forms"
	"go-user-microservice/internal/app/payment"
	"go-user-microservice/internal/app/user"
	"go-user-microservice/internal/app/user/forms/builders"
	"go-user-microservice/internal/app/wallet"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/internal/pkg/domain/providers"
)

type GrpcServerProvider struct {
	wallet  *wallet.GrpcWalletServer
	user    *user.GrpcUserServer
	card    *card.GrpcServerCard
	payment *payment.GrpcServerPayment
}

func NewGrpcServerProvider(
	serviceProvider providers.ServiceProvider,
	transactionHandler *db.TransactionHandlerDB,
) *GrpcServerProvider {
	walletGrpcServer := wallet.NewWalletGrpcServer(
		serviceProvider.WalletService(),
		transactionHandler,
	)
	userGrpcServer := user.NewUserGrpcServer(
		&builders.UserFormBuilder{},
		serviceProvider.AuthService(),
	)
	cardGrpcServer := card.NewGrpcServerCard(
		&forms.CardFormBuilder{},
		serviceProvider.CardService(),
	)
	paymentGrpcServer := payment.NewGrpcPaymentServer(serviceProvider.PaymentService(), transactionHandler)
	return &GrpcServerProvider{
		wallet:  walletGrpcServer,
		user:    userGrpcServer,
		card:    cardGrpcServer,
		payment: paymentGrpcServer,
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

func (p *GrpcServerProvider) PaymentGrpcServer() *payment.GrpcServerPayment {
	return p.payment
}
