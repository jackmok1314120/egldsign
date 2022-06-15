package egld

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Client is the interface we expect to call in order to do the HTTP requests
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	httpUserAgentKey = "User-Agent"
	httpUserAgent    = "Elrond go SDK / 1.0.0 <Posting to nodes>"

	httpAcceptTypeKey = "Accept"
	httpAcceptType    = "application/json"

	httpContentTypeKey = "Content-Type"
	httpContentType    = "application/json"
)

type ClientWrapper struct {
	url    string
	client Client
}

// NewHttpClientWrapper will create a new instance of type httpClientWrapper
func NewHttpClientWrapper(client Client, url string) *ClientWrapper {
	providedClient := client
	if IfNilReflect(providedClient) {
		providedClient = http.DefaultClient
	}

	return &ClientWrapper{
		url:    url,
		client: providedClient,
	}
}

// GetHTTP does a GET method operation on the specified endpoint
func (wrapper *ClientWrapper) GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", wrapper.url, endpoint)
	fmt.Println("url:", url)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	applyGetHeaderParams(request)

	response, err := wrapper.client.Do(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return body, response.StatusCode, nil
}

// PostHTTP does a POST method operation on the specified endpoint with the provided raw data bytes
func (wrapper *ClientWrapper) PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", wrapper.url, endpoint)
	fmt.Println("url:", url)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	applyPostHeaderParams(request)

	response, err := wrapper.client.Do(request)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	buff, err := ioutil.ReadAll(response.Body)

	return buff, response.StatusCode, err
}

// IsInterfaceNil returns true if there is no value under the interface
func (wrapper *ClientWrapper) IsInterfaceNil() bool {
	return wrapper == nil
}

func applyGetHeaderParams(request *http.Request) {
	request.Header.Set(httpAcceptTypeKey, httpAcceptType)
	request.Header.Set(httpUserAgentKey, httpUserAgent)
}

func applyPostHeaderParams(request *http.Request) {
	applyGetHeaderParams(request)
	request.Header.Set(httpContentTypeKey, httpContentType)
}

// NilInterfaceChecker checks if an interface's underlying object is nil
type NilInterfaceChecker interface {
	IsInterfaceNil() bool
}

// IfNil tests if the provided interface pointer or underlying object is nil
func IfNil(checker NilInterfaceChecker) bool {
	if checker == nil {
		return true
	}
	return checker.IsInterfaceNil()
}

// IfNilReflect tests if the provided interface pointer or underlying pointer receiver is nil
func IfNilReflect(i interface{}) bool {
	if v := reflect.ValueOf(i); v.IsValid() {
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return true
		}
		return false
	}
	return true
}
