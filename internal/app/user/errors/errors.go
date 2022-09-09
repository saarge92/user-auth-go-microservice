package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UserInnAlreadyExist  = "user inn already exists"
	UserNotFoundOnRemote = "user not found on remote server"
	SignInFail           = "password or login is incorrect"
	UserNotFound         = "user not found"
)

var (
	ErrUserAlreadyExists = status.Error(codes.AlreadyExists, UserInnAlreadyExist)
	ErrRemoteInnNotFound = status.Error(codes.NotFound, UserNotFoundOnRemote)
	ErrSignInFail        = status.Error(codes.Unauthenticated, SignInFail)
	ErrUserNotFound      = status.Error(codes.NotFound, UserNotFound)
)
