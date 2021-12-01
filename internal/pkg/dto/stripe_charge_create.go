package dto

import "github.com/shopspring/decimal"

type StripeChargeCreate struct {
	Amount   decimal.Decimal
	Currency string
}
