package rediscloud_api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

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

		mockedResponse := mockedResponses[0]
		if !mockedResponse.matches(r) {
			w.WriteHeader(501)
			return
		}

		response := mockedResponse.response()
		mockedResponses = mockedResponses[1:]
		w.WriteHeader(mockedResponse.status)
		_, _ = w.Write([]byte(response))
	}
}

type endpointRequest struct {
	method      string
	path        string
	query       url.Values
	requestBody *string
	body        string
	status      int
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
		status: http.StatusOK,
		t:      t,
	}
}

func getRequestWithQuery(t *testing.T, path string, query url.Values, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodGet,
		path:   path,
		body:   body,
		query:  query,
		status: http.StatusOK,
		t:      t,
	}
}

func getRequestWithQueryAndStatus(t *testing.T, path string, query url.Values, status int, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodGet,
		path:   path,
		body:   body,
		query:  query,
		status: status,
		t:      t,
	}
}

func getRequestWithStatus(t *testing.T, path string, status int, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodGet,
		path:   path,
		body:   body,
		status: status,
		t:      t,
	}
}

func deleteRequest(t *testing.T, path string, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodDelete,
		path:   path,
		body:   body,
		status: http.StatusOK,
		t:      t,
	}
}

func deleteRequestWithStatus(t *testing.T, path string, status int, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodDelete,
		path:   path,
		body:   body,
		status: status,
		t:      t,
	}
}

func postRequest(t *testing.T, path string, request string, body string) endpointRequest {
	return endpointRequest{
		method:      http.MethodPost,
		path:        path,
		body:        body,
		requestBody: &request,
		status:      http.StatusOK,
		t:           t,
	}
}

func postRequestWithNoRequest(t *testing.T, path string, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodPost,
		path:   path,
		body:   body,
		status: http.StatusOK,
		t:      t,
	}
}

func putRequest(t *testing.T, path string, request string, body string) endpointRequest {
	return endpointRequest{
		method:      http.MethodPut,
		path:        path,
		body:        body,
		requestBody: &request,
		status:      http.StatusOK,
		t:           t,
	}
}

func putRequestWithNoRequest(t *testing.T, path string, body string) endpointRequest {
	return endpointRequest{
		method: http.MethodPut,
		path:   path,
		body:   body,
		status: http.StatusOK,
		t:      t,
	}
}

// taskFlow returns the two endpointRequests needed for a "POST/PUT/DELETE -> GET /tasks/{id}" flow
func taskFlow(t *testing.T, method, path, requestBody, taskID, commandType string) []endpointRequest {
	now := time.Now().UTC().Format(time.RFC3339) // e.g. "2025-08-11T14:33:21Z"

	var first endpointRequest
	responseTemplate := fmt.Sprintf(`{
      "taskId": "%s",
      "commandType": "%s",
      "status": "received",
      "description": "Task queued.",
      "timestamp": "%s",
      "_links": { "task": { "href": "https://example.org", "title": "getTaskStatusUpdates", "type": "GET" } }
    }`, taskID, commandType, now)

	switch method {
	case http.MethodPost:
		if requestBody != "" {
			first = postRequest(t, path, requestBody, responseTemplate)
		} else {
			first = postRequestWithNoRequest(t, path, responseTemplate)
		}
	case http.MethodPut:
		first = putRequest(t, path, requestBody, responseTemplate)
	case http.MethodDelete:
		first = deleteRequest(t, path, responseTemplate)
	}

	completeTemplate := fmt.Sprintf(`{
      "taskId": "%s",
      "commandType": "%s",
      "status": "processing-completed",
      "timestamp": "%s",
      "response": {},
      "_links": { "self": { "href": "https://example.com", "type": "GET" } }
    }`, taskID, commandType, now)

	second := getRequest(t, "/tasks/"+taskID, completeTemplate)

	return []endpointRequest{first, second}
}
