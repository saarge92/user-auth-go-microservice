package errorlists

const (
	UserNotFoundOnRemote         = "user not found on remote server"
	RemoteServerBadAuthorization = "remote server bad request"
	UserNotFound                 = "user not found"
	TokenInvalid                 = "token is invalid"
	UserInnAlreadyExist          = "user inn already exists"
	NoInnDataRemote              = "inn data not found"
	CurrencyNotFound             = "currency not found"
	CountryNotFound              = "country not found"
	UserUnAuthenticated          = "user is not authenticated"
	UserWalletAlreadyExist       = "user wallet already exist"
	SignInFail                   = "password or login is incorrect"
	ConvertError                 = "convert error for %s"
	MustBeMore                   = "%s should be more %d"
	MustBeLess                   = "%s should be less %d"
	CardNotFound                 = "card not found"
	WalletNotFound               = "wallet not found"
)
