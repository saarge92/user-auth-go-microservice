package errorlists

const (
	RemoteServerBadAuthorization = "remote server bad request"
	TokenInvalid                 = "token is invalid"
	NoInnDataRemote              = "inn data not found"
	CurrencyNotFound             = "currency not found"
	CountryNotFound              = "country not found"
	UserUnAuthenticated          = "user is not authenticated"
	UserWalletAlreadyExist       = "user wallet already exist"
	ConvertError                 = "convert error for %s"
	MustBeMore                   = "%s should be more %d"
	MustBeLess                   = "%s should be less %d"
	WalletNotFound               = "wallet not found"
)
