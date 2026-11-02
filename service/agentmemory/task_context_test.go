package agentmemory

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMutationTaskErrorContext(t *testing.T) {
	for _, operation := range []string{"update", "delete"} {
		t.Run(operation, func(t *testing.T) {
			failure := errors.New("fixture validation failure")
			client := &fakeHTTPClient{
				patchFunc: func(ctx context.Context, name, path string, requestBody, responseBody interface{}) error {
					setTaskID(t, responseBody, "task-1")
					return nil
				},
				deleteFunc: func(ctx context.Context, name, path string, requestBody, responseBody interface{}) error {
					setTaskID(t, responseBody, "task-1")
					return nil
				},
			}
			waiter := &fakeTaskWaiter{waitFunc: func(ctx context.Context, id string) error { return failure }}
			subject := NewAPI(client, waiter, &fakeLogger{})
			var err error
			if operation == "update" {
				err = subject.Update(context.Background(), "store-1", UpdateStore{})
			} else {
				err = subject.Delete(context.Background(), "store-1")
			}
			require.ErrorIs(t, err, failure)
			require.ErrorContains(t, err, "store-1")
			if operation == "update" {
				require.ErrorContains(t, err, "updating")
			} else {
				require.ErrorContains(t, err, "deleting")
			}
		})
	}
}
