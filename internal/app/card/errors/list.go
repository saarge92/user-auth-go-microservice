package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notFoundCard = "card not found"
)

var (
	ErrCardNotFound = status.Error(codes.NotFound, notFoundCard)
)
