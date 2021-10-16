package services

import (
	entites2 "go-user-microservice/internal/app/entites"
	forms2 "go-user-microservice/internal/app/forms"
)

type WalletServiceInterface interface {
	Create(form *forms2.WalletCreateForm) (*entites2.Wallet, error)
}
