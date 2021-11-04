package dto

type StripeAccountCreate struct {
	Country               string
	Email                 string
	CardPaymentsRequested bool
	TransferRequested     bool
}
