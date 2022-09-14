package domain

import (
	"context"
	"github.com/shopspring/decimal"
	cardEntities "go-user-microservice/internal/app/card/entities"
	"go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/app/payment/filter"
	walletDto "go-user-microservice/internal/app/wallet/dto"
)

type OperationStoryRepository interface {
	Create(ctx context.Context, record *entities.OperationStory) error
	List(ctx context.Context, queryFilter *filter.OperationStoryFilter) ([]dto.OperationStory, int64, error)
}

type CardRepository interface {
	OneByCardAndUserID(ctx context.Context, externalID string, userID uint64) (*cardEntities.Card, error)
}

type WalletRepository interface {
	OneByExternalIDAndUserID(ctx context.Context, externalID string, userID uint64) (*walletDto.WalletCurrencyDto, error)
	IncreaseBalanceByID(ctx context.Context, id uint64, amount decimal.Decimal) error
}
