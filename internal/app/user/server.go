package user

import (
	"context"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/user/forms"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcUserServer struct {
	authService *services.Auth
}

func NewUserGrpcServer(
	authService *services.Auth,
) *GrpcUserServer {
	return &GrpcUserServer{
		authService: authService,
	}
}

func (s *GrpcUserServer) Signup(
	_ context.Context,
	request *core.SignUpMessage,
) (*core.SignUpResponse, error) {
	formRequest := &forms.SignUp{SignUpMessage: request}
	if e := formRequest.Validate(); e != nil {
		return nil, e
	}
	channelResponse := make(chan interface{})
	var userResponse *entities.User
	var tokenResponse string
	var errorResponse error
	go func() {
		userResponse, tokenResponse, errorResponse = s.authService.SignUp(formRequest, channelResponse)
	}()
	<-channelResponse
	if errorResponse != nil {
		return nil, errorResponse
	}
	return &core.SignUpResponse{
		Id:    userResponse.ID,
		Token: tokenResponse,
	}, nil
}

func (s *GrpcUserServer) VerifyToken(
	_ context.Context,
	request *core.VerifyMessage,
) (*core.VerifyMessageResponse, error) {
	userEntity, e := s.authService.VerifyAndReturnPayloadToken(request.Token)
	if e != nil {
		return nil, e
	}
	return &core.VerifyMessageResponse{
		User: &core.UserMessageResponse{
			Login: userEntity.Login,
			Id:    userEntity.ID,
			Roles: nil,
		},
	}, nil
}

func (s *GrpcUserServer) SignIn(
	_ context.Context,
	request *core.SignInMessage,
) (*core.SignInResponse, error) {
	formRequest := &forms.SignIn{SignInMessage: request}
	signInChan := make(chan interface{})
	var signInError error
	var userResponse *entities.User
	var token string
	go func() {
		userResponse, token, signInError = s.authService.SignIn(formRequest, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	return &core.SignInResponse{
		Id:    userResponse.ID,
		Token: token,
	}, nil
}
