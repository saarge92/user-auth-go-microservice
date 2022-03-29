package entities

import "go-user-microservice/pkg/protobuf/core"

type TransactionType uint32
type OperationType int8

const (
	DepositTransactionType  TransactionType = 1
	TransferTransactionType TransactionType = 2
	RefundTransactionType   TransactionType = 3

	DepositOperationType OperationType = 1
	RefundOperationType  OperationType = 2
)

var GRPCToOperationType = map[core.OperationType]OperationType{
	core.OperationType_Deposit: DepositOperationType,
	core.OperationType_Refund:  RefundOperationType,
}
