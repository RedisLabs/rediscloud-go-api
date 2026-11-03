package internal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpClient_Get_failsFor4xx(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	}))

	subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
	require.NoError(t, err)

	err = subject.Get(context.TODO(), "testing", "/", nil)
	require.Error(t, err)
}

func TestHttpClient_Patch(t *testing.T) {
	type requestBody struct {
		Name string `json:"name"`
	}
	type responseBody struct {
		ID string `json:"id"`
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/resource", r.URL.Path)

		var request requestBody
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		assert.Equal(t, "updated", request.Name)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"id":"resource-id"}`))
		require.NoError(t, err)
	}))

	subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
	require.NoError(t, err)

	var response responseBody
	err = subject.Patch(context.TODO(), "patch resource", "/resource", requestBody{Name: "updated"}, &response)
	require.NoError(t, err)
	assert.Equal(t, "resource-id", response.ID)
}

func TestHttpClient_allowsEmptySuccessResponse(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))

	subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
	require.NoError(t, err)

	err = subject.Delete(context.TODO(), "delete resource", "/resource", nil, nil)
	require.NoError(t, err)
}

func TestHttpClient_Retry(t *testing.T) {
	testCase := []struct {
		description   string
		retryEnabled  bool
		statusCode    int
		expectedCount int
		expectedError string
	}{
		{
			description:   "should retry 429 requests when retry is enabled",
			retryEnabled:  true,
			statusCode:    429,
			expectedCount: 3,
		},
		{
			description:   "should not retry other status code when retry is enabled",
			retryEnabled:  true,
			statusCode:    404,
			expectedCount: 1,
			expectedError: "failed to test get request: 404 - ",
		},
		{
			description:   "should not retry 429 requests when retry is disabled",
			retryEnabled:  false,
			statusCode:    429,
			expectedCount: 1,
			expectedError: "failed to test get request: 429 - ",
		},
	}

	for _, test := range testCase {
		t.Run(test.description, func(t *testing.T) {

			count := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count++
				if count < 3 {
					w.WriteHeader(test.statusCode)
					return
				}
				w.WriteHeader(200)
				_, err := w.Write([]byte("{}"))
				require.NoError(t, err)
			}))

			subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
			require.NoError(t, err)
			subject.retryEnabled = test.retryEnabled

			ctx := context.Background()
			err = subject.Get(ctx, "test get request", "/", nil)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expectedCount, count)
		})
	}

}

func TestHttpClient_RetryTransientErrors(t *testing.T) {
	closeConnection := func(t *testing.T, w http.ResponseWriter) {
		hijacker, ok := w.(http.Hijacker)
		require.True(t, ok)
		conn, _, err := hijacker.Hijack()
		require.NoError(t, err)
		require.NoError(t, conn.Close())
	}

	testCase := []struct {
		description   string
		method        string
		fail          func(t *testing.T, w http.ResponseWriter)
		expectedCount int32
		expectedError string
	}{
		{
			description:   "should retry GET on 503",
			method:        http.MethodGet,
			fail:          func(t *testing.T, w http.ResponseWriter) { w.WriteHeader(http.StatusServiceUnavailable) },
			expectedCount: 3,
		},
		{
			description:   "should retry GET on Cloudflare 521",
			method:        http.MethodGet,
			fail:          func(t *testing.T, w http.ResponseWriter) { w.WriteHeader(521) },
			expectedCount: 3,
		},
		{
			description:   "should retry GET when the server closes the connection",
			method:        http.MethodGet,
			fail:          closeConnection,
			expectedCount: 3,
		},
		{
			description:   "should not retry GET on 500",
			method:        http.MethodGet,
			fail:          func(t *testing.T, w http.ResponseWriter) { w.WriteHeader(http.StatusInternalServerError) },
			expectedCount: 1,
			expectedError: "failed to test request: 500 - ",
		},
		{
			description:   "should not retry POST on 503",
			method:        http.MethodPost,
			fail:          func(t *testing.T, w http.ResponseWriter) { w.WriteHeader(http.StatusServiceUnavailable) },
			expectedCount: 1,
			expectedError: "failed to test request: 503 - ",
		},
		{
			description:   "should not retry POST when the server closes the connection",
			method:        http.MethodPost,
			fail:          closeConnection,
			expectedCount: 1,
		},
	}

	for _, test := range testCase {
		t.Run(test.description, func(t *testing.T) {
			var count atomic.Int32
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.method, r.Method)
				if count.Add(1) < 3 {
					test.fail(t, w)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("{}"))
				require.NoError(t, err)
			}))
			defer s.Close()

			subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
			require.NoError(t, err)
			subject.retryDelay = time.Millisecond

			ctx := context.Background()
			if test.method == http.MethodGet {
				err = subject.Get(ctx, "test request", "/", nil)
			} else {
				err = subject.Post(ctx, "test request", "/", nil, nil)
			}
			switch {
			case test.expectedError != "":
				assert.EqualError(t, err, test.expectedError)
			case test.expectedCount > 1:
				assert.NoError(t, err)
			default:
				assert.Error(t, err)
			}
			assert.Equal(t, test.expectedCount, count.Load())
		})
	}
}

func TestHttpClient_Get_doesNotRetryDecodeErrors(t *testing.T) {
	var count atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
	require.NoError(t, err)
	subject.retryDelay = time.Millisecond

	var destination struct {
		ID string `json:"id"`
	}
	err = subject.Get(context.Background(), "test request", "/", &destination)
	require.ErrorIs(t, err, io.EOF)
	assert.Equal(t, int32(1), count.Load())
}

func TestHttpClient_RetryStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		cancel()
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer s.Close()

	subject, err := NewHttpClient(s.Client(), s.URL, &testLogger{t: t})
	require.NoError(t, err)
	subject.retryDelay = time.Millisecond

	err = subject.Get(ctx, "test request", "/", nil)
	require.Error(t, err)
	assert.Equal(t, int32(1), count.Load())
}

type testLogger struct {
	t *testing.T
}

func (l *testLogger) Println(v ...interface{}) {
	l.t.Log(v...)
}

var _ Log = &testLogger{}
