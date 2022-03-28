package forms

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go-user-microservice/pkg/protobuf/core"
)

type WalletCreateForm struct {
	*core.CreateWalletRequest
	UserID uint64
}

func NewWalletCreateForm(
	message *core.CreateWalletRequest,
) *WalletCreateForm {
	return &WalletCreateForm{
		CreateWalletRequest: message,
	}
}

func (f *WalletCreateForm) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.Code, is.CurrencyCode),
	)
}
