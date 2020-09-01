/*
 * ----- Authors -----
 * marcelobezer
 * cahe7cb
 * washingt0
 */

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

// HTTPClient is a shared client for http requests
type HTTPClient struct {
	// Clt Internal HTTP client
	Clt *http.Client
	// BaseURL base address for this service
	BaseURL *url.URL
	// Name defines a human readable name for this service
	Name string
}

// HTTPRequest is an http request that
// can be performed by an HTTPClient
type HTTPRequest struct {
	Inner *http.Request
	Tag   string
}

// HTTPError define errors caused by
// outgoing HTTP requests
type HTTPError struct {
	// URL the requested URL
	URL url.URL
	// Method the method used for the request
	Method string
	// StatusCode the HTTP status code returned by the remote server
	StatusCode int
	// Cause an internal error that caused the request to fail
	Cause error
	// RequestTag identifier tag of the request that caused the error
	RequestTag string
	// ClientName name of the client that attempted to perform the request
	ClientName string
}

// Error implement error interface
func (e *HTTPError) Error() string {
	return fmt.Sprintf("Request failed to [%v] '%v' with status %v: %+v", e.Method, e.URL.String(), e.StatusCode, e.Cause)
}

// Unwrap implement wrapped error interface
func (e *HTTPError) Unwrap() error {
	return e.Cause
}

// WithContext constructs a copy of the request using a context
func (req *HTTPRequest) WithContext(ctx context.Context) *HTTPRequest {
	return &HTTPRequest{Inner: req.Inner.WithContext(ctx)}
}

// WithTimeout returns a shallow copy of the request with a cancellation timeout mechanism
func (req *HTTPRequest) WithTimeout(seconds int64) (*HTTPRequest, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Inner.Context(), time.Duration(seconds)*time.Second)
	return req.WithContext(ctx), cancel
}

// WithTimeout sets a maximum timeout for outgoing requests
func (hc *HTTPClient) WithTimeout(seconds int64) *HTTPClient {
	hc.Clt.Timeout = time.Duration(seconds) * time.Second
	return hc
}

// WithName sets the readable name for this service
func (hc *HTTPClient) WithName(name string) *HTTPClient {
	hc.Name = name
	return hc
}

// NewHTTPClient creates a new HTTPClient
func NewHTTPClient(serviceBaseURL string) *HTTPClient {
	u, _ := url.Parse(serviceBaseURL)

	return &HTTPClient{
		BaseURL: u,
		Clt: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (hc *HTTPClient) buildRequest(tag string, method string, uri string, headers map[string][]string, body io.Reader, cookies ...*http.Cookie) (*HTTPRequest, error) {
	url, err := hc.BaseURL.Parse(path.Join(hc.BaseURL.Path, uri))
	if err != nil {
		return nil, err
	}

	inner, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	inner.Header = headers

	for _, cookie := range cookies {
		inner.AddCookie(cookie)
	}

	return &HTTPRequest{Inner: inner, Tag: tag}, err
}

// BuildGetRequest creates a new GET request
// to be sent using a client
func (hc *HTTPClient) BuildGetRequest(tag string, uri string, headers map[string][]string, cookies ...*http.Cookie) (*HTTPRequest, error) {
	return hc.buildRequest(tag, http.MethodGet, uri, headers, nil, cookies...)
}

// BuildPostRequest creates a new POST request
// with a payload body to be sent using a client
func (hc *HTTPClient) BuildPostRequest(tag string, uri string, headers map[string][]string, body []byte, cookies ...*http.Cookie) (*HTTPRequest, error) {
	bodyReader := ioutil.NopCloser(bytes.NewReader(body))
	return hc.buildRequest(tag, http.MethodPost, uri, headers, bodyReader, cookies...)
}

// BuildPutRequest creates a new PUT request
// with a payload body to be sent using a client
func (hc *HTTPClient) BuildPutRequest(tag string, uri string, headers map[string][]string, body []byte, cookies ...*http.Cookie) (*HTTPRequest, error) {
	bodyReader := ioutil.NopCloser(bytes.NewReader(body))
	return hc.buildRequest(tag, http.MethodPut, uri, headers, bodyReader, cookies...)
}

// BuildDeleteRequest creates a new DELETE request
// for deleting a remote resource
func (hc *HTTPClient) BuildDeleteRequest(tag string, uri string, headers map[string][]string, cookies ...*http.Cookie) (*HTTPRequest, error) {
	return hc.buildRequest(tag, http.MethodDelete, uri, headers, nil, cookies...)
}

// PerformRequest performs an HTTP request
// through a client with its environment
func (hc *HTTPClient) PerformRequest(req *HTTPRequest, data interface{}, errData error) (err error) {
	var resp *http.Response
	resp, err = hc.Clt.Do(req.Inner)
	if err != nil {
		return &HTTPError{
			URL:        *req.Inner.URL,
			Method:     req.Inner.Method,
			Cause:      err,
			StatusCode: -1,
			RequestTag: req.Tag,
			ClientName: hc.Name,
		}
	}

	defer func() {
		deferErr := resp.Body.Close()
		if deferErr != nil {
			log.Println(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNoContent:
	case http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusPartialContent:
		err = json.NewDecoder(resp.Body).Decode(data)
	default:
		if err = json.NewDecoder(resp.Body).Decode(errData); err == nil {
			err = errData
		}
	}

	if err != nil {
		return &HTTPError{
			URL:        *req.Inner.URL,
			Method:     req.Inner.Method,
			StatusCode: resp.StatusCode,
			Cause:      err,
			RequestTag: req.Tag,
			ClientName: hc.Name,
		}
	}

	return
}

// Get performs a get request and tries to unmarshall the reponse body into data
func (hc *HTTPClient) Get(tag string, uri string, headers map[string][]string, data interface{}, errData error, cookies ...*http.Cookie) error {
	req, err := hc.BuildGetRequest(tag, uri, headers, cookies...)
	if err != nil {
		return err
	}
	return hc.PerformRequest(req, data, errData)
}

// Post performs a post request and tries to unmarshall the response body into data
func (hc *HTTPClient) Post(tag string, uri string, headers map[string][]string, body []byte, data interface{}, errData error, cookies ...*http.Cookie) error {
	req, err := hc.BuildPostRequest(tag, uri, headers, body, cookies...)
	if err != nil {
		return err
	}
	return hc.PerformRequest(req, data, errData)
}

// Put performs a put request and tries to unmarshall the response body into data
func (hc *HTTPClient) Put(tag string, uri string, headers map[string][]string, body []byte, data interface{}, errData error, cookies ...*http.Cookie) error {
	req, err := hc.BuildPutRequest(tag, uri, headers, body, cookies...)
	if err != nil {
		return err
	}
	return hc.PerformRequest(req, data, errData)
}

// Delete performs a delete request of some resource and tries to unmarshall the response body into data
func (hc *HTTPClient) Delete(tag string, uri string, headers map[string][]string, data interface{}, errData error, cookies ...*http.Cookie) error {
	req, err := hc.BuildDeleteRequest(tag, uri, headers, cookies...)
	if err != nil {
		return err
	}
	return hc.PerformRequest(req, data, errData)
}

// FormatPath formats a URL Path and returns a string.
// Advantages of use this function:
// 		not need to worry about format of path parts (like "/")
// 		if this path is invalid (like if it has an invalid character) it will returns an error
func FormatPath(s ...string) (string, error) {
	if len(s) == 0 {
		return "", nil
	}

	u := &url.URL{Path: ""}
	t := &url.URL{Path: ""}
	var err error

	for i := len(s) - 1; i >= 0; i-- {
		t, err = t.Parse(s[i])
		if err != nil {
			return "", err
		}

		u, err = u.Parse(t.String() + u.String())
		if err != nil {
			return "", err
		}
	}

	return u.String(), nil
}
