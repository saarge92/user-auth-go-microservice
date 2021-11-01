package entites

type Country struct {
	ID           uint64 `db:"id"`
	CodeTwo      string `db:"code_2"`
	CodeThree    string `db:"code_3"`
	CurrencyCode string `db:"currency_code"`
	PhoneCode    int64  `db:"phone_code"`
}
