package domain

import (
	"context"
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/app/wallet/dto"
	"go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/pkg/entites"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *entities.Wallet) error
	Exist(ctx context.Context, userID uint64, currencyID uint32) (bool, error)
	DefaultWalletByUser(ctx context.Context, userID uint64) (*entities.Wallet, error)
	SetAsDefaultForUserWallet(ctx context.Context, userID uint64, isDefault bool) error
	ListByUserID(ctx context.Context, userID uint64) ([]dto.WalletCurrencyDto, error)
	OneByExternalIDAndUserID(ctx context.Context, externalID string, userID uint64) (*dto.WalletCurrencyDto, error)
	IncreaseBalanceByID(ctx context.Context, id uint64, amount decimal.Decimal) error
}

type CurrencyRepository interface {
	GetByCode(ctx context.Context, code string) (*entites.Currency, error)
}
