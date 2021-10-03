package entites

import (
	"database/sql"
	"time"
)

const (
	InnLength int = 10
)

type User struct {
	ID        uint64       `db:"id"`
	Login     string       `db:"login"`
	Inn       uint64       `db:"inn"`
	Name      string       `db:"name"`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	IsBanned  bool         `db:"is_banned"`
}
