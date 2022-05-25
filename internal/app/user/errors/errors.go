package errors

import (
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UserAlreadyExistErr  = status.Error(codes.AlreadyExists, errorlists.UserInnAlreadyExist)
	RemoteInnNotFoundErr = status.Error(codes.NotFound, errorlists.UserNotFoundOnRemote)
)
