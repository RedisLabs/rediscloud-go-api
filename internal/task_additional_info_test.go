package internal

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskPreservesAdditionalInfo(t *testing.T) {
	for _, info := range []*string{nil, func() *string { v := ""; return &v }(), func() *string { v := "/errors/invalid-data: customMemoryTypes: limit exceeded"; return &v }()} {
		name := "absent"
		if info != nil {
			name = *info
		}
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{"taskId": "task-1", "status": "processing-error", "description": "See error information for failure details.", "response": map[string]any{"error": map[string]any{"type": "AGENT_MEMORY_STORE_INVALID_REQUEST", "status": "400 BAD_REQUEST", "description": "The memory store request is invalid."}, "additionalInfo": info}})
			}))
			defer server.Close()
			client, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
			require.NoError(t, err)
			api := NewAPI(client, &testLogger{t: t})
			for _, wait := range []func() error{
				func() error { return api.Wait(context.Background(), "task-1") },
				func() error { _, err := api.WaitForResourceId(context.Background(), "task-1"); return err },
				func() error { var value any; return api.WaitForResource(context.Background(), "task-1", &value) },
			} {
				err = wait()
				require.Error(t, err)
				var taskErr *Error
				require.True(t, errors.As(err, &taskErr), "structured errors must remain discoverable")
				assert.Equal(t, "400", taskErr.StatusCode())
				assert.Contains(t, err.Error(), "task-1")
				assert.Contains(t, err.Error(), "AGENT_MEMORY_STORE_INVALID_REQUEST")
				if info != nil && *info != "" {
					assert.Contains(t, err.Error(), *info)
				} else {
					assert.Equal(t, "task task-1 failed: 400 BAD_REQUEST - AGENT_MEMORY_STORE_INVALID_REQUEST: The memory store request is invalid.", err.Error())
				}
			}
		})
	}
}

func TestHTTPEmptyResponseRequiresNoDestination(t *testing.T) {
	for _, payload := range []string{"", " \n\t"} {
		t.Run(payload, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte(payload)) }))
			defer server.Close()
			client, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
			require.NoError(t, err)
			var destination struct {
				ID string `json:"id"`
			}
			require.ErrorIs(t, client.Get(context.Background(), "read fixture", "/fixture", &destination), io.EOF)
			if payload == "" {
				require.NoError(t, client.Delete(context.Background(), "delete fixture", "/fixture", nil, nil))
			}
		})
	}
}
