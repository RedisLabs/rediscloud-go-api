package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/tags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Pro/Active-Active

// When you first create a database, the tags block is mostly empty
func TestGetTagsNonExistent(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/115251/databases/51068559/tags",
				`{
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/subscriptions/115251/databases/51068559/tags",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Tags.Get(context.TODO(), 115251, 51068559)
	require.NoError(t, err)

	assert.Equal(t, &tags.AllTags{}, actual)
}

func TestGetTagsEmpty(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/115251/databases/51068559/tags",
				`{
				  "tags": [],
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/subscriptions/115251/databases/51068559/tags",
					  "type": "GET",
					  "rel": "self"
					}
				  ],
				  "accountId": 69369
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Tags.Get(context.TODO(), 115251, 51068559)
	require.NoError(t, err)

	assert.Equal(t, &tags.AllTags{
		Tags: &[]*tags.Tag{},
	}, actual)
}

func TestGetTags(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/115251/databases/51068559/tags",
				`{
				  "tags": [
					{
					  "links": [],
					  "key": "environment",
					  "value": "production",
					  "createdAt": "2024-08-28T08:59:31.966Z",
					  "updatedAt": "2024-08-28T08:59:31.966Z"
					},
					{
					  "links": [],
					  "key": "department",
					  "value": "finance",
					  "createdAt": "2024-08-28T08:59:31.966Z",
					  "updatedAt": "2024-08-28T08:59:31.966Z"
					}
				  ],
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/subscriptions/115251/databases/51068559/tags",
					  "type": "GET",
					  "rel": "self"
					}
				  ],
				  "accountId": 69369
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Tags.Get(context.TODO(), 115251, 51068559)
	require.NoError(t, err)

	assert.Equal(t, &tags.AllTags{
		Tags: &[]*tags.Tag{
			{
				Key:   redis.String("environment"),
				Value: redis.String("production"),
			},
			{
				Key:   redis.String("department"),
				Value: redis.String("finance"),
			},
		},
	}, actual)
}

func TestPutTags(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/115251/databases/51068559/tags",
				`{
				  "tags": [
					{
					  "key": "environment",
					  "value": "production"
					},
					{
					  "key": "department",
					  "value": "finance"
					}
				  ]
				}`,
				`{
				  "tags": [
					{
					  "links": [],
					  "key": "environment",
					  "value": "production",
					  "createdAt": "2024-08-28T08:59:31.966Z",
					  "updatedAt": "2024-08-28T08:59:31.966Z"
					},
					{
					  "links": [],
					  "key": "department",
					  "value": "finance",
					  "createdAt": "2024-08-28T08:59:31.966Z",
					  "updatedAt": "2024-08-28T08:59:31.966Z"
					}
				  ],
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/subscriptions/115251/databases/51068559/tags",
					  "type": "PUT",
					  "rel": "self"
					}
				  ],
				  "accountId": 69369
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Tags.Put(
		context.TODO(),
		115251,
		51068559,
		tags.AllTags{
			Tags: &[]*tags.Tag{
				{
					Key:   redis.String("environment"),
					Value: redis.String("production"),
				},
				{
					Key:   redis.String("department"),
					Value: redis.String("finance"),
				},
			},
		},
	)

	require.NoError(t, err)
}

// Fixed

// When you first create a database, the tags block is mostly empty
func TestFixedGetTagsNonExistent(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/115229/databases/51068525/tags",
				`{
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/115229/databases/51068525/tags",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Tags.GetFixed(context.TODO(), 115229, 51068525)
	require.NoError(t, err)

	assert.Equal(t, &tags.AllTags{}, actual)
}

func TestFixedGetTagsEmpty(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/115229/databases/51068525/tags",
				`{
				  "tags": [],
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/115229/databases/51068525/tags",
					  "type": "GET",
					  "rel": "self"
					}
				  ],
				  "accountId": 69369
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Tags.GetFixed(context.TODO(), 115229, 51068525)
	require.NoError(t, err)

	assert.Equal(t, &tags.AllTags{
		Tags: &[]*tags.Tag{},
	}, actual)
}

func TestFixedGetTags(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/115229/databases/51068525/tags",
				`{
				  "tags": [
					{
					  "links": [],
					  "key": "environment",
					  "value": "production",
					  "createdAt": "2024-08-27T10:34:19.395Z",
					  "updatedAt": "2024-08-27T10:34:19.395Z"
					},
					{
					  "links": [],
					  "key": "costcenter",
					  "value": "0700",
					  "createdAt": "2024-08-27T10:34:19.395Z",
					  "updatedAt": "2024-08-27T10:34:19.395Z"
					}
				  ],
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/115229/databases/51068525/tags",
					  "type": "GET",
					  "rel": "self"
					}
				  ],
				  "accountId": 69369
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Tags.GetFixed(context.TODO(), 115229, 51068525)
	require.NoError(t, err)

	assert.Equal(t, &tags.AllTags{
		Tags: &[]*tags.Tag{
			{
				Key:   redis.String("environment"),
				Value: redis.String("production"),
			},
			{
				Key:   redis.String("costcenter"),
				Value: redis.String("0700"),
			},
		},
	}, actual)
}

func TestFixedPutTags(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/fixed/subscriptions/115229/databases/51068525/tags",
				`{
				  "tags": [
					{
					  "key": "environment",
					  "value": "production"
					},
					{
					  "key": "costCenter",
					  "value": "0700"
					}
				  ]
				}`,
				`{
				  "tags": [
					{
					  "links": [],
					  "key": "environment",
					  "value": "production",
					  "createdAt": "2024-08-27T10:34:19.395Z",
					  "updatedAt": "2024-08-27T10:34:19.395Z"
					},
					{
					  "links": [],
					  "key": "costcenter",
					  "value": "0700",
					  "createdAt": "2024-08-27T10:34:19.395Z",
					  "updatedAt": "2024-08-27T10:34:19.395Z"
					}
				  ],
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/115229/databases/51068525/tags",
					  "type": "PUT",
					  "rel": "self"
					}
				  ],
				  "accountId": 69369
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Tags.PutFixed(
		context.TODO(),
		115229,
		51068525,
		tags.AllTags{
			Tags: &[]*tags.Tag{
				{
					Key:   redis.String("environment"),
					Value: redis.String("production"),
				},
				{
					Key:   redis.String("costCenter"),
					Value: redis.String("0700"),
				},
			},
		},
	)

	require.NoError(t, err)
}
