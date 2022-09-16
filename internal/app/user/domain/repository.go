package domain

import (
	"context"
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/entites"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	UserByInnOrLoginExist(ctx context.Context, login string, inn string) (bool, error)
	UserExist(ctx context.Context, user string) (bool, error)
	GetUserWithRoles(ctx context.Context, user string) (*dto.UserRole, error)
}

type CountryRepository interface {
	GetByCodeTwo(ctx context.Context, code string) (*entites.Country, error)
}

type RoleRepository interface {
	AddUserToRole(ctx context.Context, userID uint64, roleID entities.RoleID) error
}
