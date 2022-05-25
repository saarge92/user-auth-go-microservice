package services

import (
	"bytes"
	"encoding/json"
	sharedConfig "go-user-microservice/internal/pkg/config"
	sharedErrors "go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RemoteUser struct {
	config     *sharedConfig.Config
	httpClient *http.Client
}

func NewRemoteUserService(config *sharedConfig.Config) *RemoteUser {
	return &RemoteUser{
		config:     config,
		httpClient: &http.Client{},
	}
}

func (s *RemoteUser) CheckRemoteUser(inn uint64) (r bool, e error) {
	baseURL := s.config.RemoteUserURL + "/4_1/rs/findById/party"
	token := "Token " + s.config.AuthUserRemoteKey
	postBody, e := json.Marshal(map[string]interface{}{
		"query": inn,
	})
	if e != nil {
		return false, e
	}
	bufferPostBody := bytes.NewBuffer(postBody)
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
			bodyErrorResponse := response.Body.Close()
			if bodyErrorResponse != nil && e != nil {
				e = bodyErrorResponse
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
		errorResponse := status.Error(codes.PermissionDenied, sharedErrors.RemoteServerBadAuthorization)
		return false, errorResponse
	}
	responseBody, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return false, e
	}
	var responseMap map[string]interface{}
	if e = json.Unmarshal(responseBody, &responseMap); e != nil {
		return false, e
	}
	if isVerified := s.verifyResponse(responseMap); !isVerified {
		return false, status.Error(codes.NotFound, sharedErrors.NoInnDataRemote)
	}
	return true, nil
}

func (s *RemoteUser) verifyResponse(responseMap map[string]interface{}) bool {
	if _, ok := responseMap["suggestions"]; !ok {
		return false
	}
	if sizeResponse := len(responseMap["suggestions"].([]interface{})); sizeResponse == 0 {
		return false
	}
	return true
}
