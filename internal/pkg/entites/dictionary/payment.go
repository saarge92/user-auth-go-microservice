package dictionary

type PaymentTransactionType string

const (
	DepositPaymentTransactionType    PaymentTransactionType = "deposit"
	WithdrawalPaymentTransactionType PaymentTransactionType = "withdrawal"
	TransferPaymentTransactionType   PaymentTransactionType = "transfer"
)
