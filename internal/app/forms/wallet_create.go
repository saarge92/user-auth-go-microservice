package forms

import "go-user-microservice/pkg/protobuf/wallet"

type WalletCreateForm struct {
	*wallet.CreateWalletMessage
	UserID uint64
}

func NewWalletCreateForm(
	message *wallet.CreateWalletMessage,
	userID uint64,
) *WalletCreateForm {
	return &WalletCreateForm{
		message,
		userID,
	}
}

func (f *WalletCreateForm) Validate() error {
	return nil
}
