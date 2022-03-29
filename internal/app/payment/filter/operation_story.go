package filter

import (
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/internal/pkg/filter"
)

type OperationStoryFilter struct {
	UserID        uint64
	OperationType *entities.OperationType
	filter.Pagination
}
