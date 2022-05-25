package errorlists

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CurrencyNotFoundErr = status.Error(codes.NotFound, CurrencyNotFound)
	CountryNotFoundErr  = status.Error(codes.NotFound, CountryNotFound)
)
