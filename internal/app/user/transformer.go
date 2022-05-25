package user

import (
	"go-user-microservice/internal/app/user/dto"
	"go-user-microservice/pkg/protobuf/core"
)

func FromUserResponseToGRPC(userResponse *dto.UserRole, token string) *core.SignInResponse {
	// var grpcRoles = make([]*core.Role, 0, len(userResponse.Roles))
	// for _, role := range userResponse.Roles {
	//	grpcRoles = append(grpcRoles, &core.Role{
	//		Id:   role.ID,
	//		Name: role.Name,
	//	})
	//}
	return &core.SignInResponse{
		Id:    userResponse.User.ID,
		Token: token,
		Roles: nil,
	}
}
