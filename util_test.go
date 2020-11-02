package rediscloud_api

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(apiKey, secretKey string, mockedResponses ...endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != apiKey {
			w.WriteHeader(502)
			return
		}
		if r.Header.Get("X-Api-Secret-Key") != secretKey {
			w.WriteHeader(503)
			return
		}

		if !mockedResponses[0].matches(r) {
			w.WriteHeader(501)
			return
		}

		response := mockedResponses[0].response()
		mockedResponses = mockedResponses[1:]
		w.WriteHeader(200)
		_, _ = w.Write([]byte(response))
	}
}

type endpoint interface {
	matches(r *http.Request) bool
	response() string
}

type endpointWithoutRequest struct {
	method string
	path   string
	body   string
	t      *testing.T
}

func (e endpointWithoutRequest) matches(r *http.Request) bool {
	if !assert.Equal(e.t, e.method, r.Method) {
		return false
	}
	if !assert.Equal(e.t, e.path, r.URL.Path) {
		return false
	}
	return true
}

func (e endpointWithoutRequest) response() string {
	return e.body
}

func getRequest(t *testing.T, path string, body string) endpoint {
	return endpointWithoutRequest{
		method: http.MethodGet,
		path:   path,
		body:   body,
		t:      t,
	}
}

func deleteRequest(t *testing.T, path string, body string) endpoint {
	return endpointWithoutRequest{
		method: http.MethodDelete,
		path:   path,
		body:   body,
		t:      t,
	}
}

type endpointWithRequest struct {
	method      string
	path        string
	requestBody string
	body        string
	t           *testing.T
}

func (e endpointWithRequest) matches(r *http.Request) bool {
	request, err := ioutil.ReadAll(r.Body)
	require.NoError(e.t, err)
	if !assert.JSONEq(e.t, e.requestBody, string(request)) {
		return false
	}
	if !assert.Equal(e.t, e.method, r.Method) {
		return false
	}
	if !assert.Equal(e.t, e.path, r.URL.Path) {
		return false
	}
	return true
}

func (e endpointWithRequest) response() string {
	return e.body
}

func postRequest(t *testing.T, path string, request string, body string) endpoint {
	return endpointWithRequest{
		method:      http.MethodPost,
		path:        path,
		body:        body,
		requestBody: request,
		t:           t,
	}
}

func putRequest(t *testing.T, path string, request string, body string) endpoint {
	return endpointWithRequest{
		method:      http.MethodPut,
		path:        path,
		body:        body,
		requestBody: request,
		t:           t,
	}
}
