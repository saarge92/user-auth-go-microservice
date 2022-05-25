package dto

import "go-user-microservice/internal/app/user/entities"

type UserRole struct {
	User  entities.User
	Roles []entities.Role
}
