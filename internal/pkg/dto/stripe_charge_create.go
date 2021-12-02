package dto

import "github.com/shopspring/decimal"

type StripeCardChargeCreate struct {
	Amount      decimal.Decimal
	Currency    string
	Token       string
	Description string
}
