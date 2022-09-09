package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/request"
	sharedConfig "go-user-microservice/internal/pkg/config"
	sharedErrors "go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"strings"
)

type RemoteUser struct {
	config     *sharedConfig.Config
	httpClient domain.HTTPClient
}

func NewRemoteUserService(config *sharedConfig.Config) *RemoteUser {
	return &RemoteUser{
		config:     config,
		httpClient: &http.Client{},
	}
}

func prepareRequest(ctx context.Context, method string, serviceURL string, token string, request domain.RequestParams) (*http.Request, error) {
	var body []byte
	if method == http.MethodPost {
		bodyJSON, e := json.Marshal(request)
		if e != nil {
			return nil, fmt.Errorf("make post body %w", e)
		}
		body = bodyJSON
	}

	baseURL := strings.TrimRight(serviceURL, "/") + request.RequestURI()
	preparedRequest, e := http.NewRequestWithContext(ctx, method, baseURL, bytes.NewReader(body))
	if e != nil {
		return nil, fmt.Errorf("prepare request %w", e)
	}
	contentType := "application/json"
	preparedRequest.Header = http.Header{
		"Authorization": []string{token},
		"Accept":        []string{contentType},
		"Content-Type":  []string{contentType},
	}
	return preparedRequest, nil
}

func doRequest[T domain.RequestParams](
	ctx context.Context,
	client domain.HTTPClient,
	method string,
	serviceURL string,
	token string,
	request T,
) (map[string]interface{}, error) {
	preparedRequest, e := prepareRequest(ctx, method, serviceURL, token, request)
	if e != nil {
		return nil, e
	}
	response, e := client.Do(preparedRequest)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		errorResponse := status.Error(codes.PermissionDenied, sharedErrors.RemoteServerBadAuthorization)
		return nil, errorResponse
	}

	responseBody, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	var responseMap map[string]interface{}
	if e = json.Unmarshal(responseBody, &responseMap); e != nil {
		return nil, e
	}

	return responseMap, nil
}

func (s *RemoteUser) CheckRemoteUser(ctx context.Context, request request.InnRequest) (bool, error) {
	token := "Token " + s.config.AuthUserRemoteKey
	responseMap, e := doRequest(ctx, s.httpClient, http.MethodPost, s.config.RemoteUserURL, token, request)
	if e != nil {
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
