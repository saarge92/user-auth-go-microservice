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
	userID uint64,
) *WalletCreateForm {
	return &WalletCreateForm{
		message,
		userID,
	}
}

func (f *WalletCreateForm) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.Code, is.CurrencyCode),
		validation.Field(&f.UserID, validation.Required),
	)
}
