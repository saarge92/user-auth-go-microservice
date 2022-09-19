package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	notFoundCard = "card not found"
	alreadyExist = "card already exists"
)

var (
	ErrCardNotFound     = status.Error(codes.NotFound, notFoundCard)
	ErrCardAlreadyExist = status.Error(codes.AlreadyExists, alreadyExist)
)
