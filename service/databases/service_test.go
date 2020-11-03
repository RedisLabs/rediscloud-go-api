package databases

import (
	"context"
	"fmt"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListDatabase_stopsOn404(t *testing.T) {
	client := &mockHttpClient{}

	client.On("Get", context.TODO(), "list databases for 5", "/subscriptions/5/databases?limit=100&offset=0", mock.AnythingOfType("*databases.listDatabaseResponse")).Run(func(args mock.Arguments) {
		response := args.Get(3).(*listDatabaseResponse)
		response.Subscription = []*listDbSubscription{
			{
				ID: redis.Int(5),
				Databases: []*Database{
					{
						ID: redis.Int(1),
					},
				},
			},
		}
	}).Return(nil)
	client.On("Get", context.TODO(), "list databases for 5", "/subscriptions/5/databases?limit=100&offset=100", mock.AnythingOfType("*databases.listDatabaseResponse")).Run(func(args mock.Arguments) {
		response := args.Get(3).(*listDatabaseResponse)
		response.Subscription = []*listDbSubscription{
			{
				ID: redis.Int(5),
				Databases: []*Database{
					{
						ID: redis.Int(2),
					},
				},
			},
		}
	}).Return(nil)
	client.On("Get", context.TODO(), "list databases for 5", "/subscriptions/5/databases?limit=100&offset=200", mock.AnythingOfType("*databases.listDatabaseResponse")).
		Return(&internal.HttpError{StatusCode: 404})

	subject := newListDatabase(context.TODO(), client, 5, 100)
	assert.True(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Equal(t, []*Database{
		{
			ID: redis.Int(1),
		},
	}, subject.Value())

	assert.True(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Equal(t, []*Database{
		{
			ID: redis.Int(2),
		},
	}, subject.Value())

	assert.False(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Nil(t, subject.Value())
}

func TestListDatabase_recordsError(t *testing.T) {
	client := &mockHttpClient{}

	expected := fmt.Errorf("stop")
	client.On("Get", context.TODO(), "list databases for 5", "/subscriptions/5/databases?limit=1&offset=0", mock.AnythingOfType("*databases.listDatabaseResponse")).
		Return(expected)

	subject := newListDatabase(context.TODO(), client, 5, 1)
	assert.False(t, subject.Next())
	assert.Equal(t, expected, subject.Err())
	assert.Nil(t, subject.Value())
}

type mockHttpClient struct {
	mock.Mock
}

func (m *mockHttpClient) Get(ctx context.Context, name, path string, responseBody interface{}) error {
	args := m.Called(ctx, name, path, responseBody)
	return args.Error(0)
}

func (m *mockHttpClient) Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	args := m.Called(ctx, name, path, requestBody, responseBody)
	return args.Error(0)
}

func (m *mockHttpClient) Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	args := m.Called(ctx, name, path, requestBody, responseBody)
	return args.Error(0)
}

func (m *mockHttpClient) Delete(ctx context.Context, name, path string, responseBody interface{}) error {
	args := m.Called(ctx, name, path, responseBody)
	return args.Error(0)
}
