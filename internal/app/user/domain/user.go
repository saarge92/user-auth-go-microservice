package domain

import (
	"go-user-microservice/internal/app/user/entities"
)

type UserRepositoryInterface interface {
	Create(user *entities.User) error
	Update(user *entities.User) error
	UserExist(user string) (bool, error)
	GetUser(user string) (*entities.User, error)
	UserByInnExist(inn uint64) (bool, error)
	UserByID(id uint64) (*entities.User, error)
}
