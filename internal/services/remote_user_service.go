package services

import (
	"bytes"
	"encoding/json"
	"go-user-microservice/internal/config"
	"io"
	"io/ioutil"
	"net/http"
)

type RemoteUserService struct {
	config     *config.Config
	httpClient *http.Client
}

func NewRemoteUserService(config *config.Config) *RemoteUserService {
	return &RemoteUserService{
		config:     config,
		httpClient: &http.Client{},
	}
}

func (s *RemoteUserService) CheckRemoteUser(inn uint64) (bool, error) {
	baseUrl := s.config.RemoteUserURL + "/findById/party"
	token := "Token " + s.config.AuthUserRemoteKey
	postBody, e := json.Marshal(map[string]interface{}{
		"query": inn,
	})
	bufferPostBody := bytes.NewBuffer(postBody)
	if e != nil {
		return false, e
	}
	request := &http.Request{
		Header: http.Header{
			"Authorization": []string{token},
			"Accept":        []string{"application/json"},
		},
		RequestURI: baseUrl,
		Body:       ioutil.NopCloser(bufferPostBody),
	}
	response, e := s.httpClient.Do(request)
	if e != nil {
		return false, e
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	responseBody, e := ioutil.ReadAll(response.Body)
	var responseMap map[string]interface{}
	e = json.Unmarshal(responseBody, &responseMap)
	if e != nil {
		return false, e
	}
	return true, nil
}
