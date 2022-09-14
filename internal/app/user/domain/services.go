package domain

import (
	"context"
	"github.com/stripe/stripe-go/v72"
	"go-user-microservice/internal/app/user/request"
	"go-user-microservice/internal/pkg/dto"
	"net/http"
)

type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type UserRemoteService interface {
	CheckRemoteUser(ctx context.Context, request request.InnRequest) (bool, error)
}

type StripeAccountService interface {
	Create(data *dto.StripeAccountCreate) (*stripe.Account, *stripe.Customer, error)
}
