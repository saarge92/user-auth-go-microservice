package repositories

import (
	"go-user-microservice/internal/app/entites"
)

type WalletRepositoryInterface interface {
	Create(wallet *entites.Wallet) error
}
