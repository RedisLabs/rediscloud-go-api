package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

type testLogger struct {
	t *testing.T
}

func (l *testLogger) Println(v ...interface{}) {
	l.t.Log(v...)
}

var _ Log = &testLogger{}
