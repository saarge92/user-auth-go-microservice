package repositories

import (
	"go-user-microservice/internal/entites"
)

type UserRepository interface {
	Create(user *entites.User) error
	Update(user *entites.User) error
	UserExist(user string) (bool, error)
	GetUser(user string) (*entites.User, error)
	UserByInnExist(inn uint32) (bool, error)
}
