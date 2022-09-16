package user

import (
	"context"
	"github.com/samber/lo"
	"go-user-microservice/internal/app/user/entities"
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

func (s *GrpcUserServer) VerifyToken(ctx context.Context, request *core.VerifyMessage) (*core.VerifyMessageResponse, error) {
	userRoleDto, e := s.authService.VerifyAndReturnPayloadToken(ctx, request.Token)
	if e != nil {
		return nil, e
	}
	roles := lo.Map(userRoleDto.Roles, func(role entities.Role, _ int) string {
		return role.Name
	})
	return &core.VerifyMessageResponse{
		User: &core.UserMessageResponse{
			Login: userRoleDto.User.Login,
			Id:    userRoleDto.User.ID,
			Roles: roles,
		},
	}, nil
}

func (s *GrpcUserServer) SignIn(ctx context.Context, request *core.SignInMessage) (*core.SignInResponse, error) {
	formRequest := &forms.SignIn{SignInMessage: request}
	userResponse, token, signInError := s.authService.SignIn(ctx, formRequest)
	if signInError != nil {
		return nil, signInError
	}

	return fromUserResponseToGRPC(userResponse, token), nil
}
