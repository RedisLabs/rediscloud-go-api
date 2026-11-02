package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskResourceValidationScope(t *testing.T) {
	for _, operation := range []string{"wait-without-response", "wait-without-resource", "numeric-id", "resource-object"} {
		t.Run(operation, func(t *testing.T) {
			payload := `{"status":"processing-completed"}`
			switch operation {
			case "wait-without-resource", "numeric-id":
				payload = `{"status":"processing-completed","response":{"resourceId":123}}`
			case "resource-object":
				payload = `{"status":"processing-completed","response":{"resource":{"name":"created"}}}`
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/tasks/task-1" {
					t.Errorf("unexpected task request: %s %s", r.Method, r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(payload))
			}))
			defer server.Close()
			client, err := NewHttpClient(server.Client(), server.URL, &testLogger{t: t})
			if err != nil {
				t.Fatal(err)
			}
			subject := NewAPI(client, &testLogger{t: t})
			switch operation {
			case "numeric-id":
				id, err := subject.WaitForResourceId(context.Background(), "task-1")
				if err != nil || id != 123 {
					t.Errorf("valid numeric task result rejected: %v", err)
				}
			case "resource-object":
				var value struct{ Name string }
				if err := subject.WaitForResource(context.Background(), "task-1", &value); err != nil || value.Name != "created" {
					t.Errorf("valid generic task resource rejected: %v", err)
				}
			default:
				if err := subject.Wait(context.Background(), "task-1"); err != nil {
					t.Errorf("plain task wait unexpectedly required a resource: %v", err)
				}
			}
		})
	}
}
