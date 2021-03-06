package domain

import (
	"context"
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	UserByInnOrLoginExist(ctx context.Context, login string, inn uint64) (bool, error)
	UserExist(ctx context.Context, user string) (bool, error)
	GetUserWithRoles(ctx context.Context, user string) (*dto.UserRole, error)
}
