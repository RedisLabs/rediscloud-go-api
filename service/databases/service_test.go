package databases

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListDatabase_stopsOn404(t *testing.T) {
	client := &mockHttpClient{}
	subject := newListDatabase(context.TODO(), client, 5, 100)

	client.On("GetWithQuery", context.TODO(), "list databases for 5", "/subscriptions/5/databases", url.Values{"limit": {"100"}, "offset": {"0"}}, mock.AnythingOfType("*databases.listDatabaseResponse")).Run(func(args mock.Arguments) {
		response := args.Get(4).(*listDatabaseResponse)
		response.Subscription = []*listDbSubscription{
			{
				ID: redis.Int(5),
				Databases: []*Database{
					{
						ID: redis.Int(1),
					},
					{
						ID: redis.Int(2),
					},
				},
			},
		}
	}).Return(nil)

	assert.True(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Equal(t, &Database{
		ID: redis.Int(1),
	}, subject.Value())
	assert.True(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Equal(t, &Database{
		ID: redis.Int(2),
	}, subject.Value())

	client.On("GetWithQuery", context.TODO(), "list databases for 5", "/subscriptions/5/databases", url.Values{"limit": {"100"}, "offset": {"100"}}, mock.AnythingOfType("*databases.listDatabaseResponse")).Run(func(args mock.Arguments) {
		response := args.Get(4).(*listDatabaseResponse)
		response.Subscription = []*listDbSubscription{
			{
				ID: redis.Int(5),
				Databases: []*Database{
					{
						ID: redis.Int(3),
					},
				},
			},
		}
	}).Return(nil)

	assert.True(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Equal(t, &Database{
		ID: redis.Int(3),
	}, subject.Value())

	client.On("GetWithQuery", context.TODO(), "list databases for 5", "/subscriptions/5/databases", url.Values{"limit": {"100"}, "offset": {"200"}}, mock.AnythingOfType("*databases.listDatabaseResponse")).
		Return(&internal.HTTPError{StatusCode: 404})

	assert.False(t, subject.Next())
	assert.NoError(t, subject.Err())
	assert.Nil(t, subject.Value())
}

func TestListDatabase_recordsError(t *testing.T) {
	client := &mockHttpClient{}

	expected := fmt.Errorf("stop")
	client.On("GetWithQuery", context.TODO(), "list databases for 5", "/subscriptions/5/databases", url.Values{"limit": {"1"}, "offset": {"0"}}, mock.AnythingOfType("*databases.listDatabaseResponse")).
		Return(expected)

	subject := newListDatabase(context.TODO(), client, 5, 1)
	assert.False(t, subject.Next())
	assert.Equal(t, expected, subject.Err())
	assert.Nil(t, subject.Value())
}

func TestGetCertificate_success(t *testing.T) {
	client := &mockHttpClient{}
	api := NewAPI(client, nil, nil)

	expected := &DatabaseCertificate{
		PublicCertificatePEMString: "test-certificate",
	}

	client.On("Get", context.TODO(), "get TLS certificate for database 123 in subscription 456", "/subscriptions/456/databases/123/certificate", mock.AnythingOfType("*databases.DatabaseCertificate")).
		Run(func(args mock.Arguments) {
			cert := args.Get(3).(*DatabaseCertificate)
			cert.PublicCertificatePEMString = "test-certificate"
		}).Return(nil)

	result, err := api.GetCertificate(context.TODO(), 456, 123)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetCertificate_error(t *testing.T) {
	client := &mockHttpClient{}
	api := NewAPI(client, nil, nil)

	expected := fmt.Errorf("test error")
	client.On("Get", context.TODO(), "get TLS certificate for database 123 in subscription 456", "/subscriptions/456/databases/123/certificate", mock.AnythingOfType("*databases.DatabaseCertificate")).
		Return(expected)

	result, err := api.GetCertificate(context.TODO(), 456, 123)
	assert.Error(t, err)
	assert.Equal(t, expected, err)
	assert.Nil(t, result)
}

type mockHttpClient struct {
	mock.Mock
}

func (m *mockHttpClient) Get(ctx context.Context, name, path string, responseBody interface{}) error {
	args := m.Called(ctx, name, path, responseBody)
	return args.Error(0)
}

func (m *mockHttpClient) GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error {
	args := m.Called(ctx, name, path, query, responseBody)
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

func (m *mockHttpClient) Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	args := m.Called(ctx, name, path, responseBody)
	return args.Error(0)
}
