package services

import (
	"context"
	"go-user-microservice/internal/app/user/request"
)

type RemoteUserService interface {
	CheckRemoteUser(ctx context.Context, request request.InnRequest) (r bool, e error)
}
