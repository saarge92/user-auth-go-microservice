package transformer

import (
	"go-user-microservice/internal/app/payment/dto"
	"go-user-microservice/pkg/protobuf/core"
)

func FromOperationStoriesDtoToGRPCResponse(
	operations []dto.OperationStory,
	count int64,
) *core.ListResponse {
	response := &core.ListResponse{}
	operationsGrpcElements := make([]*core.ListResponse_OperationInfo, len(operations))
	for _, operation := range operations {
		card := operation.Card
		operationStory := operation.OperationStory

		amount, _ := operationStory.Amount.Float64()
		balanceBefore, _ := operationStory.Amount.Float64()
		balanceAfter, _ := operationStory.Amount.Float64()

		grpcResponseElement := &core.ListResponse_OperationInfo{
			ExternalId:    operationStory.ExternalID,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			CreatedAt:     operation.OperationStory.CreatedAt.Unix(),
			CardInfo: &core.ListResponse_CardInfo{
				ExternalId: card.ExternalID,
				Number:     card.Number,
				CreatedAt:  card.CreatedAt.Unix(),
				IsDefault:  card.IsDefault,
			},
		}
		operationsGrpcElements = append(operationsGrpcElements, grpcResponseElement)
	}
	response.Count = count
	response.Operations = operationsGrpcElements
	return response
}
