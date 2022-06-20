package user

import (
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/pkg/protobuf/core"
)

func FromUserResponseToGRPC(userResponse *dto.UserRole, token string) *core.SignInResponse {
	return &core.SignInResponse{
		Id:    userResponse.User.ID,
		Token: token,
		Roles: nil,
	}
}
