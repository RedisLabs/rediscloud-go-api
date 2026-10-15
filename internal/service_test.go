package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompletedTaskRejectsMissingResponseWithoutPanic(t *testing.T) {
	for _, payload := range []string{
		`{"status":"processing-completed"}`,
		`{"status":"processing-completed","response":null}`,
		`{"status":"processing-completed","response":{}}`,
		`{"status":"processing-completed","response":{"resource":null,"resourceId":null}}`,
	} {
		t.Run(payload, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(payload))
			}))
			defer server.Close()
			client, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
			require.NoError(t, err)
			api := NewAPI(client, &testLogger{t: t})
			require.NotPanics(t, func() {
				var result struct{ StoreID string }
				require.ErrorContains(t, api.WaitForResource(context.Background(), "task-1", &result), "did not contain a resource")
				id, err := api.WaitForResourceId(context.Background(), "task-1")
				require.ErrorContains(t, err, "did not contain a resource ID")
				require.Zero(t, id)
			})
		})
	}
}

func TestCompletedTaskDecodesResourceAndNumericID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"processing-completed","response":{"resourceId":123,"resource":{"storeId":"store-1"}}}`))
	}))
	defer server.Close()
	client, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
	require.NoError(t, err)
	api := NewAPI(client, &testLogger{t: t})
	var result struct {
		StoreID string `json:"storeId"`
	}
	require.NoError(t, api.WaitForResource(context.Background(), "task-1", &result))
	require.Equal(t, "store-1", result.StoreID)
	id, err := api.WaitForResourceId(context.Background(), "task-1")
	require.NoError(t, err)
	require.Equal(t, 123, id)
}

func TestPendingTaskHonorsDeadlineAndCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"processing-in-progress"}`))
	}))
	defer server.Close()
	client, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
	require.NoError(t, err)
	api := NewAPI(client, &testLogger{t: t})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	require.ErrorIs(t, api.Wait(ctx, "task-1"), context.DeadlineExceeded)
	require.Less(t, time.Since(start), time.Second)
	cancelled, stop := context.WithCancel(context.Background())
	stop()
	require.ErrorIs(t, api.Wait(cancelled, "task-1"), context.Canceled)
}
