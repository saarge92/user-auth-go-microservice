package services

import (
	entites2 "go-user-microservice/internal/app/entites"
	"go-user-microservice/internal/app/forms"
)

type WalletServiceInterface interface {
	Create(form *forms.WalletCreateForm) (*entites2.Wallet, error)
}
