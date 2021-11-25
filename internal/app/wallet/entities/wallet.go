package entities

import (
	"github.com/shopspring/decimal"
	"go-user-microservice/internal/pkg/entites"
	"time"
)

type Wallet struct {
	ID         uint64          `db:"id"`
	UserID     uint64          `db:"user_id"`
	ExternalID string          `db:"external_id"`
	CurrencyID uint32          `db:"currency_id"`
	Balance    decimal.Decimal `db:"balance"`
	IsDefault  bool            `db:"is_default"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	Currency   *entites.Currency
}
