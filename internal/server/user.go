package server

import (
	"context"
	"go-user-microservice/internal/services"
	"go-user-microservice/pkg/protobuf/user"
)

type User struct {
	UserService *services.UserService
}

func (s *User) Signup(_ context.Context, in *user.SignUpMessage) (*user.SignUpResponse, error) {
	return nil, nil
}
