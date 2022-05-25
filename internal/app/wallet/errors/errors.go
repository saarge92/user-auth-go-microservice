package errors

import (
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	WalletNotFoundErr = status.Error(codes.NotFound, errorlists.WalletNotFound)
)
