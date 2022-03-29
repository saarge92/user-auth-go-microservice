package transformers

import (
	"go-user-microservice/internal/app/payment/entities"
	"go-user-microservice/pkg/protobuf/core"
)

func FromGRPCOperationType(operation core.OperationType) *entities.OperationType {
	if operation == core.OperationType_ALL {
		return nil
	}
	result := entities.GRPCToOperationType[operation]
	return &result
}
