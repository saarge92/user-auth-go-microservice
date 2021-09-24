package repositories

import "go-user-microservice/internal/entites"

type UserRepository interface {
	Create(user *entites.User) error
	Update(user *entites.User) error
}
