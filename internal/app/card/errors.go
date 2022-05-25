package card

import (
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CardNotFoundErr = status.Error(codes.NotFound, errorlists.CardNotFound)
)
