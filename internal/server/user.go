package server

import (
	"context"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/forms"
	"go-user-microservice/internal/services"
	"go-user-microservice/pkg/protobuf/user"
)

type UserGrpcServer struct {
	userService *services.UserService
}

func NewUserGrpcServer(userService *services.UserService) *UserGrpcServer {
	return &UserGrpcServer{
		userService: userService,
	}
}

func (s *UserGrpcServer) Signup(_ context.Context, in *user.SignUpMessage) (*user.SignUpResponse, error) {
	form := &forms.SignUp{SignUpMessage: in}
	channelResponse := make(chan interface{})
	var userResponse *entites.User
	var errorResponse error
	go func() {
		userResponse, errorResponse = s.userService.SignUp(form, channelResponse)
	}()
	<-channelResponse
	if errorResponse != nil {
		return nil, errorResponse
	}
	return nil, nil
}
