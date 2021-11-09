package card

type Card struct {
	ID                 uint64 `db:"id"`
	Number             string `db:"number"`
	ExpireDay          uint32 `db:"expire_day"`
	ExpireMonth        uint32 `db:"expire_month"`
	ExpireYear         uint32 `db:"expire_year"`
	UserID             uint64 `db:"user_id"`
	ExternalProviderID string `db:"external_provider_id"`
	IsDefault          bool   `db:"is_default"`
}
