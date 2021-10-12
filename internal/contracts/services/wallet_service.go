package services

import (
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/forms"
)

type WalletServiceInterface interface {
	Create(form *forms.WalletCreateForm) (*entites.Wallet, error)
}
