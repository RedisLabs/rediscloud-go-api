package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumericTaskResultValidationScope(t *testing.T) {
	for _, tc := range []struct {
		name      string
		status    int
		payload   string
		wantError string
		plainWait bool
		cancel    bool
	}{
		{name: "explicit-zero-unchanged", status: http.StatusOK, payload: `{"status":"processing-completed","response":{"resourceId":0}}`},
		{name: "task-http-error", status: http.StatusInternalServerError, wantError: "500"},
		{name: "failed-task", status: http.StatusOK, payload: `{"status":"processing-failed","description":"simulated task failure"}`, wantError: "simulated task failure"},
		{name: "cancelled-context", cancel: true},
		{name: "plain-wait-needs-no-id", status: http.StatusOK, payload: `{"status":"processing-completed"}`, plainWait: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.payload))
			}))
			defer server.Close()
			c, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
			require.NoError(t, err)
			subject := NewAPI(c, &testLogger{t: t})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancel {
				cancel()
			}
			if tc.plainWait {
				require.NoError(t, subject.Wait(ctx, "task-1"))
				return
			}
			id, err := subject.WaitForResourceId(ctx, "task-1")
			if tc.cancel {
				require.ErrorIs(t, err, context.Canceled)
			} else if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
			} else {
				require.NoError(t, err)
			}
			assert.Zero(t, id)
		})
	}
}
