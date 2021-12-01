package form

import "go-user-microservice/pkg/protobuf/payment"

type Deposit struct {
	*payment.DepositRequest
}
