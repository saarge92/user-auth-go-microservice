package repositories

import "go-user-microservice/internal/entites"

type WalletRepositoryInterface interface {
	Create(wallet *entites.Wallet) error
}
