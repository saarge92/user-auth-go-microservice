package services

import (
	"context"
	"go-user-microservice/internal/app/wallet/entities"
	"go-user-microservice/internal/app/wallet/forms"
)

type WalletServiceInterface interface {
	Create(ctx context.Context, form *forms.WalletCreateForm) (*entities.Wallet, error)
}
