package forms

import "go-user-microservice/pkg/protobuf/wallet"

type WalletCreateForm struct {
	*wallet.CreateWalletMessage
	UserID uint64
}
