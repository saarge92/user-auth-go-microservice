package entities

type TransactionType uint32
type OperationType int8

const (
	DepositTransactionType  TransactionType = 1
	TransferTransactionType TransactionType = 2
	RefundTransactionType   TransactionType = 3

	DepositOperationType OperationType = 1
	RefundOperationType  OperationType = 2
)
