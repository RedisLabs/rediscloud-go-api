package rediscloud_api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clientFromTestServer(s *httptest.Server, apiKey string, secretKey string) (*Client, error) {
	return NewClient(LogRequests(true), BaseURL(s.URL), Auth(apiKey, secretKey), Transporter(s.Client().Transport))
}

func testServer(apiKey, secretKey string, mockedResponses ...endpointRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "go-http-client") {
			w.WriteHeader(504)
			return
		}
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

type endpointRequest struct {
	method      string
	path        string
	query       url.Values
	requestBody *string
	body        string
	t           *testing.T
}

func (e endpointRequest) matches(r *http.Request) bool {
	if e.requestBody != nil {
		request, err := ioutil.ReadAll(r.Body)
		require.NoError(e.t, err)
		if !assert.JSONEq(e.t, *e.requestBody, string(request)) {
			return false
		}
	} else {
		if !assert.Empty(e.t, r.Body) {
			return false
		}
	}

	if e.query != nil {
		if !assert.Equal(e.t, e.query, r.URL.Query()) {
			return false
		}
	} else {
		if !assert.Empty(e.t, r.URL.RawQuery) {
			return false
		}
	}

	if !assert.Equal(e.t, e.method, r.Method) {
		return false
	}
	if !assert.Equal(e.t, e.path, r.URL.Path) {
		return false
	}
	return true
}

func (e endpointRequest) response() string {
	return e.body
}

func getRequest(t *testing.T, path string, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodGet,
		path:   path,
		body:   body,
		t:      t,
	}
}

func getRequestWithQuery(t *testing.T, path string, query url.Values, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodGet,
		path:   path,
		body:   body,
		query:  query,
		t:      t,
	}
}

func deleteRequest(t *testing.T, path string, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodDelete,
		path:   path,
		body:   body,
		t:      t,
	}
}

func postRequest(t *testing.T, path string, request string, body string) endpointRequest {
	return endpointRequest{
		method:      http.MethodPost,
		path:        path,
		body:        body,
		requestBody: &request,
		t:           t,
	}
}

func postRequestWithNoRequest(t *testing.T, path string, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodPost,
		path:   path,
		body:   body,
		t:      t,
	}
}

func putRequest(t *testing.T, path string, request string, body string) endpointRequest {
	return endpointRequest{
		method:      http.MethodPut,
		path:        path,
		body:        body,
		requestBody: &request,
		t:           t,
	}
}
