package user

import (
	"github.com/samber/lo"
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/pkg/protobuf/core"
)

func fromUserResponseToGRPC(userResponse *dto.UserRole, token string) *core.SignInResponse {
	roles := lo.Map(userResponse.Roles, func(role entities.Role, _ int) string {
		return role.Name
	})
	return &core.SignInResponse{
		Id:    userResponse.User.ID,
		Token: token,
		Roles: roles,
	}
}
