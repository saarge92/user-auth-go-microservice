package services

import (
	"context"
	"go-user-microservice/internal/app/wallet/forms"
	"go-user-microservice/internal/pkg/entites"
)

type WalletServiceInterface interface {
	Create(ctx context.Context, form *forms.WalletCreateForm) (*entites.Wallet, error)
}
