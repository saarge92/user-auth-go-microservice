package repositories

import (
	"go-user-microservice/internal/entites"
)

type UserRepositoryInterface interface {
	Create(user *entites.User) error
	Update(user *entites.User) error
	UserExist(user string) (bool, error)
	GetUser(user string) (*entites.User, error)
	UserByInnExist(inn uint64) (bool, error)
	UserByID(id uint64) (*entites.User, error)
}
