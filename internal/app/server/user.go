package server

import (
	"context"
	"go-user-microservice/internal/app/entites"
	userFormBuilders "go-user-microservice/internal/app/forms/user/builders"
	"go-user-microservice/internal/app/services/member"
	protoMember "go-user-microservice/pkg/protobuf/member"
)

type UserGrpcServer struct {
	authService     *member.AuthService
	userFormBuilder *userFormBuilders.UserFormBuilder
}

func NewUserGrpcServer(
	userFormBuilder *userFormBuilders.UserFormBuilder,
	authService *member.AuthService,
) *UserGrpcServer {
	return &UserGrpcServer{
		userFormBuilder: userFormBuilder,
		authService:     authService,
	}
}

func (s *UserGrpcServer) Signup(
	_ context.Context,
	request *protoMember.SignUpMessage,
) (*protoMember.SignUpResponse, error) {
	form := s.userFormBuilder.Signup(request)
	if e := form.Validate(); e != nil {
		return nil, e
	}
	channelResponse := make(chan interface{})
	var userResponse *entites.User
	var tokenResponse string
	var errorResponse error
	go func() {
		userResponse, tokenResponse, errorResponse = s.authService.SignUp(form, channelResponse)
	}()
	<-channelResponse
	if errorResponse != nil {
		return nil, errorResponse
	}
	return &protoMember.SignUpResponse{
		Id:    userResponse.ID,
		Token: tokenResponse,
	}, nil
}

func (s *UserGrpcServer) VerifyToken(
	_ context.Context,
	request *protoMember.VerifyMessage,
) (*protoMember.VerifyMessageResponse, error) {
	userPayload, userEntity, e := s.authService.VerifyAndReturnPayloadToken(request.Token)
	if e != nil {
		return nil, e
	}
	return &protoMember.VerifyMessageResponse{
		User: &protoMember.UserMessageResponse{
			Login: userPayload.UserName,
			Id:    userEntity.ID,
			Roles: nil,
		},
	}, nil
}

func (s *UserGrpcServer) SignIn(
	_ context.Context,
	request *protoMember.SignInMessage,
) (*protoMember.SignInResponse, error) {
	form := s.userFormBuilder.SignIn(request)
	signInChan := make(chan interface{})
	var signInError error
	var userResponse *entites.User
	var token string
	go func() {
		userResponse, token, signInError = s.authService.SignIn(form, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	return &protoMember.SignInResponse{
		Id:    userResponse.ID,
		Token: token,
	}, nil
}
