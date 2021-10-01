package services

import (
	"bytes"
	"encoding/json"
	"go-user-microservice/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (s *RemoteUserService) CheckRemoteUser(inn uint32) (r bool, e error) {
	baseURL := s.config.RemoteUserURL + "/4_1/rs/findById/party"
	token := "Token " + s.config.AuthUserRemoteKey
	postBody, e := json.Marshal(map[string]interface{}{
		"query": inn,
	})
	bufferPostBody := bytes.NewBuffer(postBody)
	if e != nil {
		return false, e
	}
	baseRequestURL, e := url.Parse(baseURL)
	if e != nil {
		return false, e
	}
	request := &http.Request{
		Method: "POST",
		Header: http.Header{
			"Authorization": []string{token},
			"Accept":        []string{"application/json"},
			"Content-Type":  []string{"application/json"},
		},
		URL:  baseRequestURL,
		Body: ioutil.NopCloser(bufferPostBody),
	}
	response, e := s.httpClient.Do(request)
	defer func() {
		if response != nil {
			e = response.Body.Close()
			if e != nil {
				return
			}
		}
	}()
	if e != nil {
		return false, e
	}
	if response == nil {
		return false, nil
	}
	if response.StatusCode == http.StatusBadRequest {
		errorResponse := status.Error(codes.PermissionDenied, "")
		return false, errorResponse
	}
	responseBody, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return false, e
	}
	var responseMap map[string]interface{}
	e = json.Unmarshal(responseBody, &responseMap)
	if e != nil {
		return false, e
	}
	return true, nil
}
