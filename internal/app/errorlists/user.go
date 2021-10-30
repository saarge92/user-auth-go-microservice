package errorlists

const (
	UserNotFoundOnRemote         = "user not found on remote server"
	RemoteServerBadAuthorization = "remote server bad request"
	UserNotFound                 = "user not found"
	TokenInvalid                 = "token is invalid"
	UserEmailAlreadyExist        = "user already exist"
	UserInnAlreadyExist          = "user inn already exists"
	NoInnDataRemote              = "inn data not found"
	CurrencyNotFound             = "currency not found"
	UserUnAuthenticated          = "user is not authenticated"
	UserWalletAlreadyExist       = "user wallet already exist"
	SignInFail                   = "password or login is incorrect"
	ConvertError                 = "convert error for %s"
)
