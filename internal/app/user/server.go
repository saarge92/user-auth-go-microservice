package user

import (
	"context"
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/user/forms"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/db"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcUserServer struct {
	authService        *services.Auth
	transactionHandler *db.TransactionHandlerDB
}

func NewUserGrpcServer(
	authService *services.Auth,
	transactionHandler *db.TransactionHandlerDB,
) *GrpcUserServer {
	return &GrpcUserServer{
		authService:        authService,
		transactionHandler: transactionHandler,
	}
}

func (s *GrpcUserServer) Signup(
	ctx context.Context,
	request *core.SignUpMessage,
) (signupResponse *core.SignUpResponse, e error) {
	formRequest := &forms.SignUp{SignUpMessage: request}
	if e = formRequest.Validate(); e != nil {
		return nil, e
	}
	ctx, tx, e := s.transactionHandler.Create(ctx, nil)
	if e != nil {
		return
	}
	defer func() {
		e = db.HandleTransaction(tx, e)
	}()

	channelResponse := make(chan interface{})
	var userResponse *entities.User
	var tokenResponse string
	var errorResponse error
	go func() {
		userResponse, tokenResponse, errorResponse = s.authService.SignUp(ctx, formRequest, channelResponse)
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
	userEntity, e := s.authService.VerifyAndReturnPayloadToken(context.Background(), request.Token)
	if e != nil {
		return nil, e
	}
	return &core.VerifyMessageResponse{
		User: &core.UserMessageResponse{
			Login: userEntity.User.Login,
			Id:    userEntity.User.ID,
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
	var userResponse *dto.UserRole
	var token string
	go func() {
		userResponse, token, signInError = s.authService.SignIn(context.Background(), formRequest, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	return FromUserResponseToGRPC(userResponse, token), nil
}
