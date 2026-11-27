package rediscloud_api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/service/psc"
)

// These are simulated API shapes, not evidence of deployed Cloud behavior.
func TestPSCReadTaskValidationHTTP(t *testing.T) {
	ctx := context.Background()
	operations := []struct {
		name, path, query string
		call              func(*Client) (any, error)
	}{
		{"service", "/subscriptions/1/private-service-connect", "", func(c *Client) (any, error) { return c.PrivateServiceConnect.GetService(ctx, 1) }},
		{"service_active_active", "/subscriptions/1/regions/2/private-service-connect", "", func(c *Client) (any, error) { return c.PrivateServiceConnect.GetActiveActiveService(ctx, 1, 2) }},
		{"endpoints", "/subscriptions/1/private-service-connect/2", "", func(c *Client) (any, error) { return c.PrivateServiceConnect.GetEndpoints(ctx, 1, 2) }},
		{"endpoints_active_active", "/subscriptions/1/regions/2/private-service-connect/3", "", func(c *Client) (any, error) { return c.PrivateServiceConnect.GetActiveActiveEndpoints(ctx, 1, 2, 3) }},
		{"creation_script", "/subscriptions/1/private-service-connect/2/endpoints/3/creationScripts", "includeTerraformGcpScript=true", func(c *Client) (any, error) {
			return c.PrivateServiceConnect.GetEndpointCreationScripts(ctx, 1, 2, 3, true)
		}},
		{"creation_script_active_active", "/subscriptions/1/regions/2/private-service-connect/3/endpoints/4/creationScripts", "includeTerraformGcpScript=true", func(c *Client) (any, error) {
			return c.PrivateServiceConnect.GetActiveActiveEndpointCreationScripts(ctx, 1, 2, 3, 4, true)
		}},
		{"deletion_script", "/subscriptions/1/private-service-connect/2/endpoints/3/deletionScripts", "", func(c *Client) (any, error) { return c.PrivateServiceConnect.GetEndpointDeletionScripts(ctx, 1, 2, 3) }},
		{"deletion_script_active_active", "/subscriptions/1/regions/2/private-service-connect/3/endpoints/4/deletionScripts", "", func(c *Client) (any, error) {
			return c.PrivateServiceConnect.GetActiveActiveEndpointDeletionScripts(ctx, 1, 2, 3, 4)
		}},
	}
	cases := []struct {
		name, initial, task       string
		initialStatus, taskStatus int
		wantErr                   bool
		wantTask                  int32
	}{
		{"missing", `{}`, "", 200, 0, true, 0},
		{"null", `{"taskId":null}`, "", 200, 0, true, 0},
		{"empty", `{"taskId":""}`, "", 200, 0, true, 0},
		{"empty_body", ``, "", 200, 0, true, 0},
		{"invalid_json", `{`, "", 200, 0, true, 0},
		{"initial_forbidden", `{"message":"fixture-error"}`, "", 403, 0, true, 0},
		{"initial_not_found", `{"message":"fixture-error"}`, "", 404, 0, true, 0},
		{"task_forbidden", `{"taskId":"test-task"}`, `{"message":"fixture-error"}`, 200, 403, true, 1},
		{"missing_resource", `{"taskId":"test-task"}`, `{"status":"processing-completed","response":{}}`, 200, 200, true, 1},
		{"valid", `{"taskId":"test-task"}`, `{"status":"processing-completed","response":{"resource":{"id":7,"status":"active"}}}`, 200, 200, false, 1},
	}
	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					var initialCalls, taskCalls atomic.Int32
					srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodGet || r.Header.Get("X-Api-Key") != "fixture-key" || r.Header.Get("X-Api-Secret-Key") != "fixture-secret" {
							t.Error("unexpected method or credentials")
							http.Error(w, "unexpected", http.StatusForbidden)
							return
						}
						w.Header().Set("Content-Type", "application/json")
						switch r.URL.Path {
						case op.path:
							if r.URL.RawQuery != op.query {
								t.Errorf("unexpected query %q", r.URL.RawQuery)
							}
							initialCalls.Add(1)
							w.WriteHeader(tc.initialStatus)
							_, _ = fmt.Fprint(w, tc.initial)
						case "/tasks/test-task":
							if r.URL.RawQuery != "" {
								t.Error("unexpected task query")
							}
							taskCalls.Add(1)
							w.WriteHeader(tc.taskStatus)
							_, _ = fmt.Fprint(w, tc.task)
						default:
							t.Errorf("unexpected path %q", r.URL.Path)
							http.Error(w, "unexpected", http.StatusNotFound)
						}
					}))
					defer srv.Close()
					api, err := NewClient(BaseURL(srv.URL), Auth("fixture-key", "fixture-secret"), Transporter(srv.Client().Transport))
					if err != nil {
						t.Fatal(err)
					}
					result, err := op.call(api)
					if (err != nil) != tc.wantErr || initialCalls.Load() != 1 || taskCalls.Load() != tc.wantTask {
						t.Fatalf("unexpected error, initial calls=%d task calls=%d", initialCalls.Load(), taskCalls.Load())
					}
					if tc.name == "initial_not_found" {
						var single *psc.NotFound
						var multi *psc.NotFoundActiveActive
						if strings.HasSuffix(op.name, "_active_active") {
							if !errors.As(err, &multi) {
								t.Fatal("active-active not-found type not preserved")
							}
						} else if !errors.As(err, &single) {
							t.Fatal("not-found type not preserved")
						}
					}
					if !tc.wantErr && (result == nil || reflect.ValueOf(result).IsNil()) {
						t.Fatal("valid response not decoded")
					}
				})
			}
		})
	}
}
