package dto

import "github.com/shopspring/decimal"

type StripeCardCustomerChargeCreate struct {
	Amount      decimal.Decimal
	Currency    string
	CardID      string
	CustomerID  string
	Description string
}
