package rediscloud_api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
		Body:          ioutil.NopCloser(bytes.NewBufferString(`{"Here":"Value"}`)),
		ContentLength: 16,
		Host:          "example.org",
	}
	expected := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Response":"Value"}`)),
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
		Body:          ioutil.NopCloser(bytes.NewBufferString(`{"password":"pass"}`)),
		ContentLength: 19,
		Host:          "example.org",
	}
	expected := &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     map[string][]string{},
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"password":"REDACTED"}`)),
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
