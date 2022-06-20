package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type WalletCurrencyDto struct {
	Wallet struct {
		ID         uint64          `db:"wallet.id"`
		UserID     uint64          `db:"wallet.user_id"`
		ExternalID string          `db:"wallet.external_id"`
		CurrencyID uint32          `db:"wallet.currency_id"`
		Balance    decimal.Decimal `db:"wallet.balance"`
		IsDefault  bool            `db:"wallet.is_default"`
		CreatedAt  time.Time       `db:"wallet.created_at"`
		UpdatedAt  time.Time       `db:"wallet.updated_at"`
	}
	Currency struct {
		Code string `db:"currency.code"`
	}
}
