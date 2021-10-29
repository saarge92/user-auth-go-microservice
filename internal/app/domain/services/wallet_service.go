package services

import (
	"context"
	"go-user-microservice/internal/app/entites"
	"go-user-microservice/internal/app/forms"
)

type WalletServiceInterface interface {
	Create(ctx context.Context, form *forms.WalletCreateForm) (*entites.Wallet, error)
}
