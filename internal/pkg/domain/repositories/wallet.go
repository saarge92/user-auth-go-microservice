package repositories

import (
	"go-user-microservice/internal/pkg/entites"
)

type WalletRepositoryInterface interface {
	Create(wallet *entites.Wallet) error
	Exist(userID uint64, currencyID uint32) (bool, error)
}
