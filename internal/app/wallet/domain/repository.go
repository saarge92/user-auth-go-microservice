package domain

import (
	"context"
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/internal/app/wallet/entities"
)

type WalletRepositoryInterface interface {
	Create(ctx context.Context, wallet *entities.Wallet) error
	Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error)
	ByUserAndDefault(ctx context.Context, userID uint64, isDefault bool) (*entities.Wallet, error)
	UpdateStatusByUserID(ctx context.Context, userID uint64, isDefault bool) error
	ListByUserID(ctx context.Context, userID uint64) ([]dto.WalletCurrencyDto, error)
}
