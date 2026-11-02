package rediscloud_api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/agentmemory"
)

func TestAgentMemoryCreateTaskResource(t *testing.T) {
	for _, scenario := range []struct {
		name       string
		response   string
		postStatus int
		getStatus  int
		wantError  string
	}{
		{name: "response-omitted", wantError: "did not contain a resource"},
		{name: "response-null", response: `,"response":null`, wantError: "did not contain a resource"},
		{name: "response-empty", response: `,"response":{}`, wantError: "did not contain a resource"},
		{name: "resource-omitted", response: `,"response":{"resourceId":123}`, wantError: "did not contain a resource"},
		{name: "resource-null", response: `,"response":{"resource":null}`, wantError: "did not contain a resource"},
		{name: "valid", response: `,"response":{"resource":{"storeId":"store-1"}}`},
		{name: "empty-resource-object", response: `,"response":{"resource":{}}`, wantError: "did not contain a store ID"},
		{name: "wrong-resource-type", response: `,"response":{"resource":123}`, wantError: "cannot unmarshal"},
		{name: "task-get-error", getStatus: http.StatusInternalServerError, wantError: "retrieve Task"},
		{name: "create-post-error", postStatus: http.StatusInternalServerError, wantError: "create Agent Memory store"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			var posts, gets atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/memory-stores":
					posts.Add(1)
					if scenario.postStatus != 0 {
						w.WriteHeader(scenario.postStatus)
						return
					}
					_, _ = w.Write([]byte(`{"taskId":"task-1"}`))
				case r.Method == http.MethodGet && r.URL.Path == "/tasks/task-1":
					gets.Add(1)
					if scenario.getStatus != 0 {
						w.WriteHeader(scenario.getStatus)
						return
					}
					_, _ = w.Write([]byte(`{"status":"processing-completed"` + scenario.response + `}`))
				default:
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			subject, err := NewClient(BaseURL(server.URL), Auth("fake-key", "fake-secret"))
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			// Malformed task completion must return an error rather than crash its caller.
			defer func() {
				if recover() != nil {
					t.Error("Agent Memory Create panicked on a completed task without a resource")
				}
			}()
			id, err := subject.AgentMemory.Create(ctx, agentmemory.CreateStore{Name: redis.String("task-resource-test"), DatabaseID: redis.Int(123)})
			if scenario.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), scenario.wantError) {
					t.Errorf("expected error containing %q, got %v", scenario.wantError, err)
				}
				if id != "" {
					t.Error("failed creation returned a fabricated store ID")
				}
			} else if err != nil || id != "store-1" {
				t.Errorf("valid completion did not return the created store: %v", err)
			}
			wantGets := int32(1)
			if scenario.postStatus != 0 {
				wantGets = 0
			}
			if posts.Load() != 1 || gets.Load() != wantGets {
				t.Error("operation was retried or polled an unacknowledged task")
			}
		})
	}
}
