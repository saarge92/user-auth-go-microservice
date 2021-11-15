package dto

type StripeAccountCreate struct {
	Country      string
	Email        string
	CardPayments bool
	Transfers    bool
}
