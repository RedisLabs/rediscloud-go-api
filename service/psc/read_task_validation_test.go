package psc

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type readTaskClient struct {
	id  *string
	err error
}

func (c *readTaskClient) Get(_ context.Context, _, _ string, out interface{}) error {
	if c.err != nil {
		return c.err
	}
	out.(*internal.TaskResponse).ID = c.id
	return nil
}
func (c *readTaskClient) GetWithQuery(ctx context.Context, name, path string, _ url.Values, out interface{}) error {
	return c.Get(ctx, name, path, out)
}
func (*readTaskClient) Post(context.Context, string, string, interface{}, interface{}) error {
	panic("unexpected POST")
}
func (*readTaskClient) Put(context.Context, string, string, interface{}, interface{}) error {
	panic("unexpected PUT")
}
func (*readTaskClient) Delete(context.Context, string, string, interface{}, interface{}) error {
	panic("unexpected DELETE")
}

type readTaskWaiter struct {
	calls int
	id    string
}

func (*readTaskWaiter) WaitForResourceId(context.Context, string) (int, error) {
	panic("unexpected resource ID wait")
}
func (*readTaskWaiter) Wait(context.Context, string) error { panic("unexpected wait") }
func (w *readTaskWaiter) WaitForResource(_ context.Context, id string, out interface{}) error {
	w.calls++
	w.id = id
	return json.Unmarshal([]byte(`{}`), out)
}

type readTaskLogger struct{}

func (readTaskLogger) Printf(string, ...interface{}) {}

func TestPSCReadRejectsMissingTaskID(t *testing.T) {
	ctx := context.Background()
	empty, valid := "", "task-fixture"
	upstream := errors.New("fixture transport error")
	operations := []struct {
		name string
		call func(*API) error
	}{
		{"service", func(a *API) error { _, e := a.GetService(ctx, 1); return e }},
		{"service_active_active", func(a *API) error { _, e := a.GetActiveActiveService(ctx, 1, 2); return e }},
		{"endpoints", func(a *API) error { _, e := a.GetEndpoints(ctx, 1, 2); return e }},
		{"endpoints_active_active", func(a *API) error { _, e := a.GetActiveActiveEndpoints(ctx, 1, 2, 3); return e }},
		{"creation_script", func(a *API) error { _, e := a.GetEndpointCreationScripts(ctx, 1, 2, 3, true); return e }},
		{"creation_script_active_active", func(a *API) error { _, e := a.GetActiveActiveEndpointCreationScripts(ctx, 1, 2, 3, 4, true); return e }},
		{"deletion_script", func(a *API) error { _, e := a.GetEndpointDeletionScripts(ctx, 1, 2, 3); return e }},
		{"deletion_script_active_active", func(a *API) error { _, e := a.GetActiveActiveEndpointDeletionScripts(ctx, 1, 2, 3, 4); return e }},
	}
	for _, op := range operations {
		t.Run(op.name, func(t *testing.T) {
			for _, tc := range []struct {
				name     string
				id       *string
				upstream error
				wantErr  bool
				wantWait int
			}{
				{"missing", nil, nil, true, 0}, {"empty", &empty, nil, true, 0}, {"valid", &valid, nil, false, 1}, {"transport_error", nil, upstream, true, 0},
			} {
				t.Run(tc.name, func(t *testing.T) {
					waiter := &readTaskWaiter{}
					err := op.call(NewAPI(&readTaskClient{id: tc.id, err: tc.upstream}, waiter, readTaskLogger{}))
					if (err != nil) != tc.wantErr || waiter.calls != tc.wantWait {
						t.Fatalf("error=%v, waits=%d", err, waiter.calls)
					}
					if tc.upstream != nil && !errors.Is(err, upstream) {
						t.Fatalf("transport error not preserved: %v", err)
					}
					if tc.upstream == nil && tc.wantErr && err.Error() != "private service connect task response did not contain a task ID" {
						t.Fatalf("unexpected error: %v", err)
					}
					if tc.wantWait == 1 && waiter.id != valid {
						t.Fatalf("wrong task ID: %q", waiter.id)
					}
				})
			}
		})
	}
}
