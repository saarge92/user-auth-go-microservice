package entities

import (
	"database/sql"
)

const (
	InnLength int = 10
)

type User struct {
	ID                 uint64        `db:"id"`
	Login              string        `db:"login"`
	Inn                string        `db:"inn"`
	Name               string        `db:"name"`
	Password           string        `db:"password"`
	CreatedAt          int64         `db:"created_at"`
	UpdatedAt          int64         `db:"updated_at"`
	DeletedAt          sql.NullInt64 `db:"deleted_at"`
	IsBanned           bool          `db:"is_banned"`
	CountryID          sql.NullInt64 `db:"country_id"`
	AccountProviderID  string        `db:"account_provider_id"`
	CustomerProviderID string        `db:"customer_provider_id"`
}
