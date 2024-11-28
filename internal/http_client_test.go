package internal

import (
	"context"
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
