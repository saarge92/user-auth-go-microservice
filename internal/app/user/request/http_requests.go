package request

type InnRequest struct {
	Inn uint64 `json:"query"`
}

func (r InnRequest) RequestURI() string {
	return "/4_1/rs/findById/party"
}
