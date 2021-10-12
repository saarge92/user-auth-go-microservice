package entites

import (
	"github.com/shopspring/decimal"
	"time"
)

type Wallet struct {
	ID         uint64          `db:"id"`
	UserID     uint64          `db:"user_id"`
	CurrencyID uint32          `db:"currency_id"`
	Balance    decimal.Decimal `db:"balance"`
	IsDefault  bool            `db:"is_default"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
	Currency   *Currency
}
