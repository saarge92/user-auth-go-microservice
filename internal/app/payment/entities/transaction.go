package entities

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID                 uint64          `db:"id"`
	ExternalID         string          `db:"external_id"`
	TransactionType    TransactionType `db:"transaction_type"`
	ExternalProviderID string          `db:"external_provider_id"`
	FromUserID         uint64          `db:"from_user_id"`
	ToUserID           uint64          `db:"to_user_id"`
	Amount             decimal.Decimal `db:"amount"`
	CreatedAt          time.Time       `db:"created_at"`
}
