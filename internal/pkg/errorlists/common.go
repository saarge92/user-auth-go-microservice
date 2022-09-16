package errorlists

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCurrencyNotFound    = status.Error(codes.NotFound, CurrencyNotFound)
	ErrCountryNotFound     = status.Error(codes.NotFound, CountryNotFound)
	ErrUserUnAuthenticated = status.Error(codes.Unauthenticated, UserUnAuthenticated)
)
