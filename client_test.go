package rediscloud_api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCredentialTripper_DisabledLoggingLogsNothing(t *testing.T) {
	mockTripper := &mockedRoundTripper{}
	mockLogger := &mockedLogger{}

	request := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.org",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Host:       "example.org",
	}
	expected := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	mockTripper.On("RoundTrip", request).Return(expected, nil)

	subject := &credentialTripper{
		apiKey:      "KEY THAT SHOULD NOT BE LOGGED",
		secretKey:   "SECRET KEY THAT SHOULD NOT BE LOGGED",
		wrapped:     mockTripper,
		logRequests: false,
		logger:      mockLogger,
		userAgent:   "test-user-agent",
	}

	actual, err := subject.RoundTrip(request)
	require.NoError(t, err)
	assert.Same(t, expected, actual)
	assert.Empty(t, mockLogger.log)
	assert.Equal(t, "KEY THAT SHOULD NOT BE LOGGED", request.Header.Get("X-Api-Key"))
	assert.Equal(t, "SECRET KEY THAT SHOULD NOT BE LOGGED", request.Header.Get("X-Api-Secret-Key"))
}

func TestCredentialTripper_ModelKeysRedactedWithoutChangingPayloads(t *testing.T) {
	const requestBody = `{"llm":{"credentials":{"type":"apiKey","apiKey":"FAKE-LLM-\"quoted\"-\\slash"}},"embedding":{"credentials":{"type":"apiKey","apiKey":"FAKE-EMBEDDING"}},"security":{"password":"FAKE-PASSWORD","global_password":"FAKE-GLOBAL-PASSWORD"},"keyName":"keep-name","apiKeyId":"keep-id"}`
	const responseBody = `{"apiKey":"FAKE-GENERATED-\"quoted\"-\\slash","keyName":"keep-name","apiKeyId":"keep-id"}`
	for _, enabled := range []bool{true, false} {
		t.Run(fmt.Sprintf("logging=%t", enabled), func(t *testing.T) {
			received := make(chan string, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				received <- string(body)
				w.Header().Set("Content-Type", "application/json")
				_, err = io.WriteString(w, responseBody)
				assert.NoError(t, err)
			}))
			defer server.Close()

			logger := &mockedLogger{}
			httpClient := &http.Client{Transport: &credentialTripper{
				apiKey: "fake-access", secretKey: "fake-secret",
				wrapped: server.Client().Transport, logRequests: enabled,
				logger: logger, userAgent: "test-user-agent",
			}}
			request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, server.URL+"/memory-stores", strings.NewReader(requestBody))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			response, err := httpClient.Do(request)
			require.NoError(t, err)
			defer func() { assert.NoError(t, response.Body.Close()) }()
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			assert.Equal(t, requestBody, <-received)
			assert.Equal(t, responseBody, string(body))

			if !enabled {
				assert.Empty(t, logger.log)
				return
			}
			require.Len(t, logger.log, 2)
			_, requestJSON, found := strings.Cut(logger.log[0], "\n\n")
			require.True(t, found, "request body was not logged")
			_, responseJSON, found := strings.Cut(logger.log[1], "\n\n")
			require.True(t, found, "response body was not logged")
			assert.JSONEq(t, `{"llm":{"credentials":{"type":"apiKey","apiKey":"REDACTED"}},"embedding":{"credentials":{"type":"apiKey","apiKey":"REDACTED"}},"security":{"password":"REDACTED","global_password":"REDACTED"},"keyName":"keep-name","apiKeyId":"keep-id"}`, requestJSON)
			assert.JSONEq(t, `{"apiKey":"REDACTED","keyName":"keep-name","apiKeyId":"keep-id"}`, responseJSON)
			for _, secret := range []string{"FAKE-LLM", "FAKE-EMBEDDING", "FAKE-PASSWORD", "FAKE-GLOBAL-PASSWORD", "FAKE-GENERATED"} {
				assert.NotContains(t, strings.Join(logger.log, "\n"), secret)
			}
		})
	}
}

func TestCredentialTripper_LogsRequestAndResponseBodies(t *testing.T) {
	mockTripper := &mockedRoundTripper{}
	mockLogger := &mockedLogger{}

	request := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.org",
			Path:   "/foo/bar",
		},
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        map[string][]string{},
		Body:          io.NopCloser(bytes.NewBufferString(`{"Here":"Value"}`)),
		ContentLength: 16,
		Host:          "example.org",
	}
	expected := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Body:       io.NopCloser(bytes.NewBufferString(`{"Response":"Value"}`)),
	}

	mockTripper.On("RoundTrip", request).Return(expected, nil)

	subject := &credentialTripper{
		apiKey:      "KEY THAT SHOULD NOT BE LOGGED",
		secretKey:   "SECRET KEY THAT SHOULD NOT BE LOGGED",
		wrapped:     mockTripper,
		logRequests: true,
		logger:      mockLogger,
		userAgent:   "test-user-agent",
	}

	actual, err := subject.RoundTrip(request)
	require.NoError(t, err)
	assert.Same(t, expected, actual)
	assert.Equal(t, "KEY THAT SHOULD NOT BE LOGGED", request.Header.Get("X-Api-Key"))
	assert.Equal(t, "SECRET KEY THAT SHOULD NOT BE LOGGED", request.Header.Get("X-Api-Secret-Key"))

	assert.Equal(t, `DEBUG: Request /foo/bar:
---[ REQUEST ]---
POST /foo/bar HTTP/1.1
Host: example.org
User-Agent: test-user-agent
Content-Length: 16
Accept: application/json
Accept-Encoding: gzip

{
  "Here": "Value"
}`, mockLogger.log[0])
	assert.Equal(t, `DEBUG: Response /foo/bar:
---[ RESPONSE ]---
HTTP/1.1 200 OK
Connection: close

{
  "Response": "Value"
}`, mockLogger.log[1])
}

func TestCredentialTripper_HandlesNoBodies(t *testing.T) {
	mockTripper := &mockedRoundTripper{}
	mockLogger := &mockedLogger{}

	request := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/baz",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Host:       "example.org",
	}
	expected := &http.Response{
		StatusCode: http.StatusNoContent,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: map[string][]string{
			"X-Test": {"Demo"},
		},
	}

	mockTripper.On("RoundTrip", request).Return(expected, nil)

	subject := &credentialTripper{
		apiKey:      "KEY THAT SHOULD NOT BE LOGGED",
		secretKey:   "SECRET KEY THAT SHOULD NOT BE LOGGED",
		wrapped:     mockTripper,
		logRequests: true,
		logger:      mockLogger,
		userAgent:   "test-user-agent",
	}

	actual, err := subject.RoundTrip(request)
	require.NoError(t, err)
	assert.Same(t, expected, actual)
	assert.Equal(t, "KEY THAT SHOULD NOT BE LOGGED", request.Header.Get("X-Api-Key"))
	assert.Equal(t, "SECRET KEY THAT SHOULD NOT BE LOGGED", request.Header.Get("X-Api-Secret-Key"))

	assert.Equal(t, `DEBUG: Request /baz:
---[ REQUEST ]---
GET /baz HTTP/1.1
Host: example.org
User-Agent: test-user-agent
Accept: application/json
Accept-Encoding: gzip`, mockLogger.log[0])
	assert.Equal(t, `DEBUG: Response /baz:
---[ RESPONSE ]---
HTTP/1.1 204 No Content
X-Test: Demo`, mockLogger.log[1])
}

func TestCredentialTripper_RedactPasswordFromBody(t *testing.T) {
	mockTripper := &mockedRoundTripper{}
	mockLogger := &mockedLogger{}

	request := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.org",
			Path:   "/foo/bar",
		},
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        map[string][]string{},
		Body:          io.NopCloser(bytes.NewBufferString(`{"password":"pass"}`)),
		ContentLength: 19,
		Host:          "example.org",
	}
	expected := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Body:       io.NopCloser(bytes.NewBufferString(`{"password":"REDACTED"}`)),
	}

	mockTripper.On("RoundTrip", request).Return(expected, nil)

	subject := &credentialTripper{
		apiKey:      "KEY THAT SHOULD NOT BE LOGGED",
		secretKey:   "SECRET KEY THAT SHOULD NOT BE LOGGED",
		wrapped:     mockTripper,
		logRequests: true,
		logger:      mockLogger,
		userAgent:   "test-user-agent",
	}

	actual, err := subject.RoundTrip(request)
	require.NoError(t, err)
	assert.Same(t, expected, actual)

	assert.Equal(t, `DEBUG: Request /foo/bar:
---[ REQUEST ]---
POST /foo/bar HTTP/1.1
Host: example.org
User-Agent: test-user-agent
Content-Length: 19
Accept: application/json
Accept-Encoding: gzip

{
  "password": "REDACTED"
}`, mockLogger.log[0])
}

func TestCredentialTripper_RedactPasswordFromNestedBody(t *testing.T) {
	mockTripper := &mockedRoundTripper{}
	mockLogger := &mockedLogger{}

	request := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.org",
			Path:   "/foo/bar",
		},
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        map[string][]string{},
		Body:          io.NopCloser(bytes.NewBufferString(`{"security": {"password":"pass", "global_password":"globalpass"}}`)),
		ContentLength: 65,
		Host:          "example.org",
	}
	expected := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Body:       io.NopCloser(bytes.NewBufferString(`{"security": {"password":"REDACTED", "global_password":"REDACTED"}}`)),
	}

	mockTripper.On("RoundTrip", request).Return(expected, nil)

	subject := &credentialTripper{
		apiKey:      "KEY THAT SHOULD NOT BE LOGGED",
		secretKey:   "SECRET KEY THAT SHOULD NOT BE LOGGED",
		wrapped:     mockTripper,
		logRequests: true,
		logger:      mockLogger,
		userAgent:   "test-user-agent",
	}

	actual, err := subject.RoundTrip(request)
	require.NoError(t, err)
	assert.Same(t, expected, actual)

	assert.Equal(t, `DEBUG: Request /foo/bar:
---[ REQUEST ]---
POST /foo/bar HTTP/1.1
Host: example.org
User-Agent: test-user-agent
Content-Length: 65
Accept: application/json
Accept-Encoding: gzip

{
  "security": {
    "password": "REDACTED",
    "global_password": "REDACTED"
  }
}`, mockLogger.log[0])
}

type mockedLogger struct {
	log []string
}

func (m *mockedLogger) Printf(format string, v ...interface{}) {
	m.log = append(m.log, m.trim(fmt.Sprintf(format, v...)))
}

func (m *mockedLogger) Println(v ...interface{}) {
	for _, i := range v {
		m.log = append(m.log, m.trim(fmt.Sprintf("%s", i)))
	}
}

func (m mockedLogger) trim(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\r", ""))
}

var _ Log = &mockedLogger{}

type mockedRoundTripper struct {
	mock.Mock
}

func (m *mockedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	args := m.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}

var _ http.RoundTripper = &mockedRoundTripper{}
