package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClient_Get_failsFor4xx(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	}))

	subject, err := NewHttpClient(s.Client(), s.URL)
	require.NoError(t, err)

	err = subject.Get(context.TODO(), "testing", "/", nil)
	require.Error(t, err)
}
