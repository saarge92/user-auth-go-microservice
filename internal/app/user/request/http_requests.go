package request

type InnRequest struct {
	Inn string `json:"query"`
}

func (r InnRequest) RequestURI() string {
	return "/4_1/rs/findById/party"
}
