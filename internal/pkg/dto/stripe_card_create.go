package dto

type StripeCardCreate struct {
	Number             string
	ExpireMonth        uint8
	ExpireYear         uint32
	CVC                uint32
	AccountProviderID  string
	CustomerProviderID string
}
