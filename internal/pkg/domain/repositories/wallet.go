package repositories

import (
	"context"
	"go-user-microservice/internal/pkg/entites"
)

type WalletRepositoryInterface interface {
	Create(ctx context.Context, wallet *entites.Wallet) error
	Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error)
	ByUserAndDefault(ctx context.Context, userID uint64, isDefault bool) (*entites.Wallet, error)
	UpdateStatusByUserID(ctx context.Context, userID uint64, isDefault bool) error
}
