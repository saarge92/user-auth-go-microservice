package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go-user-microservice/pkg/protobuf/wallet"
)

type WalletCreateForm struct {
	*wallet.CreateWalletMessage
	UserID uint64
}

func NewWalletCreateForm(
	message *wallet.CreateWalletMessage,
) *WalletCreateForm {
	return &WalletCreateForm{
		CreateWalletMessage: message,
	}
}

func (f *WalletCreateForm) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.Code, is.CurrencyCode),
	)
}
