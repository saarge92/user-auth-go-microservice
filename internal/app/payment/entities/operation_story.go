package entities

import (
	"github.com/shopspring/decimal"
	"time"
)

type OperationStory struct {
	ID                 uint64          `db:"id"`
	ExternalID         string          `db:"external_id"`
	UserID             uint64          `db:"user_id"`
	CardID             uint64          `db:"card_id"`
	Amount             decimal.Decimal `db:"amount"`
	BalanceBefore      decimal.Decimal `db:"balance_before"`
	BalanceAfter       decimal.Decimal `db:"balance_after"`
	ExternalProviderID string          `db:"external_provider_id"`
	OperationTypeID    OperationType   `db:"operation_type_id"`
	CreatedAt          time.Time       `db:"created_at"`
}
