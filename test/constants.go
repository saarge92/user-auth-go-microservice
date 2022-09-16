package test

import "go-user-microservice/internal/app/user/entities"

const (
	CardNumberForCreate          = "5555555555554444"
	UserID                uint64 = 1
	UserAccountProviderID        = "account-uuid"
	UserCustomerID               = "customer-uuid"
)

var (
	CurrentUser = &entities.User{
		ID:                 UserID,
		AccountProviderID:  UserAccountProviderID,
		CustomerProviderID: UserCustomerID,
	}
)
