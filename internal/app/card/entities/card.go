package entities

import (
	"database/sql"
	"time"
)

type Card struct {
	ID                 uint64       `db:"id"`
	ExternalID         string       `db:"external_id"`
	Number             string       `db:"number"`
	ExpireMonth        uint32       `db:"expire_month"`
	ExpireYear         uint32       `db:"expire_year"`
	UserID             uint64       `db:"user_id"`
	ExternalProviderID string       `db:"external_provider_id"`
	IsDefault          bool         `db:"is_default"`
	CreatedAt          time.Time    `db:"created_at"`
	UpdatedAt          time.Time    `db:"updated_at"`
	DeleteAt           sql.NullTime `db:"updated_at"`
}
