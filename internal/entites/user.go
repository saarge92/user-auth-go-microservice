package entites

import (
	"database/sql"
	"time"
)

const (
	InnLength     int = 10
	NameLengthMin     = 2
	NameLengthMax     = 120
)

type User struct {
	ID        int32        `db:"id"`
	Login     string       `db:"login"`
	Name      string       `db:"name"`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	IsBanned  bool         `db:"is_banned"`
}
