package dto

import (
	"go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/pkg/entites"
)

type WalletCurrencyDto struct {
	entities.Wallet  `db:"wallet"`
	entites.Currency `db:"currency"`
}
