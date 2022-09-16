package test

import (
	userDto "go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/app/user/entities"
)

const (
	CardNumberForCreate          = "5555555555554444"
	UserID                uint64 = 1
	UserAccountProviderID        = "account-uuid"
	UserCustomerID               = "customer-uuid"
)

var UserRoleData = &userDto.UserRole{
	User: entities.User{
		ID:                 UserID,
		AccountProviderID:  UserAccountProviderID,
		CustomerProviderID: UserCustomerID,
	},
}
