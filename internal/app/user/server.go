package user

import (
	"context"
	"go-user-microservice/internal/app/user/forms"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/database"
	"go-user-microservice/pkg/protobuf/core"
)

type GrpcUserServer struct {
	authService        *services.Auth
	transactionHandler *database.TransactionHandlerDB
}

func NewUserGrpcServer(
	authService *services.Auth,
	transactionHandler *database.TransactionHandlerDB,
) *GrpcUserServer {
	return &GrpcUserServer{
		authService:        authService,
		transactionHandler: transactionHandler,
	}
}

func (s *GrpcUserServer) Signup(ctx context.Context, request *core.SignUpMessage) (*core.SignUpResponse, error) {
	formRequest := &forms.SignUp{SignUpMessage: request}
	if e := formRequest.Validate(); e != nil {
		return nil, e
	}

	transactionHandler := database.NewTypedTransaction[*core.SignUpResponse](s.transactionHandler)

	return transactionHandler.WithCtx(ctx, func(ctx context.Context) (*core.SignUpResponse, error) {
		userResponse, tokenResponse, errorResponse := s.authService.SignUp(ctx, formRequest)
		if errorResponse != nil {
			return nil, errorResponse
		}
		return &core.SignUpResponse{
			Id:    userResponse.ID,
			Token: tokenResponse,
		}, nil
	})
}

func (s *GrpcUserServer) VerifyToken(_ context.Context, request *core.VerifyMessage) (*core.VerifyMessageResponse, error) {
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

func (s *GrpcUserServer) SignIn(ctx context.Context, request *core.SignInMessage) (*core.SignInResponse, error) {
	formRequest := &forms.SignIn{SignInMessage: request}
	userResponse, token, signInError := s.authService.SignIn(ctx, formRequest)
	if signInError != nil {
		return nil, signInError
	}

	return FromUserResponseToGRPC(userResponse, token), nil
}
