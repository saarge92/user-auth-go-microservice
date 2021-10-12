package entites

import "time"

type Currency struct {
	ID          uint32    `db:"id"`
	Code        string    `db:"code"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
