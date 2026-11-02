package rediscloud_api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/redis_rules"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
)

func TestCreateNumericResourceRequiresTaskIDResult(t *testing.T) {
	for _, operation := range []struct {
		name string
		path string
		call func(*Client) (int, error)
	}{
		{name: "subscription", path: "/subscriptions", call: func(c *Client) (int, error) {
			return c.Subscription.Create(context.Background(), subscriptions.CreateSubscription{})
		}},
		{name: "database", path: "/subscriptions/123/databases", call: func(c *Client) (int, error) {
			return c.Database.Create(context.Background(), 123, databases.CreateDatabase{})
		}},
		{name: "acl-rule", path: "/acl/redisRules", call: func(c *Client) (int, error) {
			return c.RedisRules.Create(context.Background(), redis_rules.CreateRedisRuleRequest{})
		}},
	} {
		t.Run(operation.name, func(t *testing.T) {
			for _, tc := range []struct {
				name      string
				payload   string
				wantError bool
			}{
				{name: "response-omitted", payload: `{"status":"processing-completed"}`, wantError: true},
				{name: "response-null", payload: `{"status":"processing-completed","response":null}`, wantError: true},
				{name: "response-empty", payload: `{"status":"processing-completed","response":{}}`, wantError: true},
				{name: "id-omitted", payload: `{"status":"processing-completed","response":{"resource":{}}}`, wantError: true},
				{name: "id-null", payload: `{"status":"processing-completed","response":{"resourceId":null}}`, wantError: true},
				{name: "valid-id", payload: `{"status":"processing-completed","response":{"resourceId":456}}`},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var creates, reads atomic.Int32
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						switch {
						case r.Method == http.MethodPost && r.URL.Path == operation.path:
							creates.Add(1)
							_, _ = w.Write([]byte(`{"taskId":"task-1"}`))
						case r.Method == http.MethodGet && r.URL.Path == "/tasks/task-1":
							reads.Add(1)
							_, _ = w.Write([]byte(tc.payload))
						default:
							t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
							http.NotFound(w, r)
						}
					}))
					defer server.Close()
					c, err := NewClient(BaseURL(server.URL), Auth("fake-key", "fake-secret"))
					require.NoError(t, err)
					var id int
					require.NotPanics(t, func() { id, err = operation.call(c) })
					assert.EqualValues(t, 1, creates.Load())
					assert.EqualValues(t, 1, reads.Load())
					if tc.wantError {
						require.ErrorContains(t, err, "completed task response did not contain a resource ID")
						assert.Zero(t, id)
						return
					}
					require.NoError(t, err)
					assert.Equal(t, 456, id)
				})
			}
		})
	}
}
