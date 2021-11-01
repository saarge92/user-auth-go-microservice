package entites

import "math/big"

type PaymentTransaction struct {
	ID big.Int `db:"id"`
}
