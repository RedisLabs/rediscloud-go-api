package rediscloud_api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Create(t *testing.T) {
	expected := 4291
	s := httptest.NewServer(testServer("key", "secret", postRequest(t, "/subscriptions/42/databases", `{
  "dryRun": false,
  "name": "Redis-database-example",
  "protocol": "redis",
  "datasetSizeInGb": 1,
  "supportOSSClusterApi": false,
  "respVersion": "resp3",
  "useExternalEndpointForOSSClusterApi": false,
  "dataPersistence": "none",
  "dataEvictionPolicy": "allkeys-lru",
  "queryPerformanceFactor": "6x",
  "redisVersion": "6.0.5",
  "replication": true,
  "throughputMeasurement": {
    "by": "operations-per-second",
    "value": 1000
  },
  "averageItemSizeInBytes": 1,
  "ramPercentage": 20,
  "replicaOf": [
    "another"
  ],
  "sourceIp": [
    "10.0.0.1"
  ],
  "enableTls": true,
  "clientSslCertificate": "something",
  "clientTlsCertificates": ["something", "else"],
  "password": "fooBar",
  "alerts": [
    {
      "name": "dataset-size",
      "value": 80
    }
  ],
  "modules": [
    {
      "name": "RedisSearch"
    }
  ],
  "autoMinorVersionUpgrade": true
}`, `{
  "taskId": "task",
  "commandType": "databaseCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", fmt.Sprintf(`{
  "taskId": "task",
  "commandType": "databaseCreateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resourceId": %d
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`, expected))))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Database.Create(context.TODO(), 42, databases.CreateDatabase{
		DryRun:                              redis.Bool(false),
		Name:                                redis.String("Redis-database-example"),
		Protocol:                            redis.String("redis"),
		DatasetSizeInGB:                     redis.Float64(1),
		SupportOSSClusterAPI:                redis.Bool(false),
		RespVersion:                         redis.String("resp3"),
		UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
		DataPersistence:                     redis.String("none"),
		DataEvictionPolicy:                  redis.String("allkeys-lru"),
		QueryPerformanceFactor:              redis.String("6x"),
		RedisVersion:                        redis.String("6.0.5"),
		Replication:                         redis.Bool(true),
		ThroughputMeasurement: &databases.CreateThroughputMeasurement{
			By:    redis.String("operations-per-second"),
			Value: redis.Int(1000),
		},
		AverageItemSizeInBytes: redis.Int(1),
		RamPercentage:          redis.Int(20),
		ReplicaOf:              redis.StringSlice("another"),
		SourceIP:               redis.StringSlice("10.0.0.1"),
		ClientSSLCertificate:   redis.String("something"),
		ClientTLSCertificates:  &[]*string{redis.String("something"), redis.String("else")},
		EnableTls:              redis.Bool(true),
		Password:               redis.String("fooBar"),
		Alerts: []*databases.Alert{
			{
				Name:  redis.String("dataset-size"),
				Value: redis.Int(80),
			},
		},
		Modules: []*databases.Module{
			{
				Name: redis.String("RedisSearch"),
			},
		},
		AutoMinorVersionUpgrade: redis.Bool(true),
	})
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestDatabase_List(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequestWithQuery(t, "/subscriptions/23456/databases", map[string][]string{"limit": {"100"}, "offset": {"0"}}, `{
  "accountId": 2,
  "subscription": [
    {
      "subscriptionId": 23456,
      "databases": [
        {
          "databaseId": 42,
          "name": "first-example",
          "protocol": "redis",
          "provider": "AWS",
          "region": "eu-west-1",
		  "queryPerformanceFactor": "Standard",
          "redisVersion": "6.0.5",
		  "ramPercentage": 20
        },
        {
          "databaseId": 43,
          "name": "second-example",
          "protocol": "redis",
          "provider": "AWS",
          "region": "eu-west-1",
		  "queryPerformanceFactor": "Standard",
          "redisVersion": "6.0.5",
  		  "ramPercentage": 30
        }
      ]
    }
  ],
  "_links": {
    "self": {
      "href": "https://example.org",
      "type": "GET"
    }
  }
}`), getRequestWithQueryAndStatus(t, "/subscriptions/23456/databases", map[string][]string{"limit": {"100"}, "offset": {"100"}}, 404, "")))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	list := subject.Database.List(context.TODO(), 23456)

	var actual []*databases.Database
	for list.Next() {
		actual = append(actual, list.Value())
	}
	require.NoError(t, list.Err())

	assert.Equal(t, []*databases.Database{
		{
			ID:                     redis.Int(42),
			Name:                   redis.String("first-example"),
			Protocol:               redis.String("redis"),
			Provider:               redis.String("AWS"),
			Region:                 redis.String("eu-west-1"),
			QueryPerformanceFactor: redis.String("Standard"),
			RedisVersion:           redis.String("6.0.5"),
			RamPercentage:          redis.Int(20),
		},
		{
			ID:                     redis.Int(43),
			Name:                   redis.String("second-example"),
			Protocol:               redis.String("redis"),
			Provider:               redis.String("AWS"),
			Region:                 redis.String("eu-west-1"),
			QueryPerformanceFactor: redis.String("Standard"),
			RedisVersion:           redis.String("6.0.5"),
			RamPercentage:          redis.Int(30),
		},
	}, actual)

	assert.False(t, redis.BoolValue(actual[0].ActiveActiveRedis))
	assert.False(t, redis.BoolValue(actual[1].ActiveActiveRedis))
}

func TestDatabase_Get(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/23456/databases/98765", `{
  "databaseId": 98765,
  "name": "Example",
  "protocol": "redis",
  "provider": "AWS",
  "region": "eu-west-1",
  "redisVersionCompliance": "6.0.5",
  "status": "active",
  "memoryLimitInGb": 7,
  "datasetSizeInGb": 7,
  "ramPercentage": 20,
  "memoryUsedInMb": 5,
  "memoryStorage": "ram",
  "supportOSSClusterApi": true,
  "respVersion": "resp2",
  "dataPersistence": "none",
  "replication": false,
  "dataEvictionPolicy": "volatile-random",
  "throughputMeasurement": {
    "by": "operations-per-second",
    "value": 10000
  },
  "QueryPerformanceFactor": "Standard",
  "redisVersion": "6.0.5",
  "activatedOn": "2020-11-03T09:03:30Z",
  "lastModified": "2020-11-03T09:03:30Z",
  "publicEndpoint": "example.com:16668",
  "privateEndpoint": "example.net:16668",
  "replicaOf": {
    "endpoints": ["another"]
  },
  "clustering": {
    "numberOfShards": 1,
    "regexRules": [
      {
        "ordinal": 1,
        "pattern": "(?<tag>.*)"
      },
      {
        "ordinal": 0,
        "pattern": ".*\\{(?<tag>.*)\\}.*"
      }
    ],
    "hashingPolicy": "custom"
  },
  "security": {
	"enableDefaultUser": false,
    "password": "test",
    "sslClientAuthentication": false,
	"tlsClientAuthentication": true,
    "enableTls": true,
    "sourceIps": [
      "0.0.0.0/0"
    ]
  },
  "modules": [],
  "alerts": [],
  "_links": {
    "self": {
      "href": "https://example.org",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Database.Get(context.TODO(), 23456, 98765)
	require.NoError(t, err)

	assert.Equal(t, &databases.Database{
		ID:                   redis.Int(98765),
		Name:                 redis.String("Example"),
		Protocol:             redis.String("redis"),
		Provider:             redis.String("AWS"),
		Region:               redis.String("eu-west-1"),
		Status:               redis.String("active"),
		MemoryLimitInGB:      redis.Float64(7),
		DatasetSizeInGB:      redis.Float64(7),
		RamPercentage:        redis.Int(20),
		MemoryUsedInMB:       redis.Float64(5),
		SupportOSSClusterAPI: redis.Bool(true),
		RespVersion:          redis.String("resp2"),
		DataPersistence:      redis.String("none"),
		Replication:          redis.Bool(false),
		ReplicaOf: &databases.ReplicaOf{
			Endpoints: []*string{redis.String("another")},
		},
		DataEvictionPolicy:     redis.String("volatile-random"),
		ActivatedOn:            redis.Time(time.Date(2020, 11, 3, 9, 3, 30, 0, time.UTC)),
		LastModified:           redis.Time(time.Date(2020, 11, 3, 9, 3, 30, 0, time.UTC)),
		MemoryStorage:          redis.String("ram"),
		PrivateEndpoint:        redis.String("example.net:16668"),
		PublicEndpoint:         redis.String("example.com:16668"),
		RedisVersionCompliance: redis.String("6.0.5"),
		ThroughputMeasurement: &databases.Throughput{
			By:    redis.String("operations-per-second"),
			Value: redis.Int(10_000),
		},
		QueryPerformanceFactor: redis.String("Standard"),
		RedisVersion:           redis.String("6.0.5"),
		Clustering: &databases.Clustering{
			NumberOfShards: redis.Int(1),
			RegexRules: []*databases.RegexRule{
				{
					Ordinal: 1,
					Pattern: "(?<tag>.*)",
				},
				{
					Ordinal: 0,
					Pattern: ".*\\{(?<tag>.*)\\}.*",
				},
			},
		},
		Security: &databases.Security{
			EnableDefaultUser:       redis.Bool(false),
			SSLClientAuthentication: redis.Bool(false),
			TLSClientAuthentication: redis.Bool(true),
			EnableTls:               redis.Bool(true),
			SourceIPs:               redis.StringSlice("0.0.0.0/0"),
			Password:                redis.String("test"),
		},
		Modules: []*databases.Module{},
		Alerts:  []*databases.Alert{},
	}, actual)
}

func TestDatabase_Get_wraps404Error(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequestWithStatus(t, "/subscriptions/23456/databases/98765", 404, "")))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Database.Get(context.TODO(), 23456, 98765)

	assert.Nil(t, actual)
	assert.IsType(t, &databases.NotFound{}, err)
}

func TestDatabase_Update(t *testing.T) {
	flow := taskFlow(
		t,
		http.MethodPut,
		"/subscriptions/42/databases/18",
		`{
		  "dryRun": false,
		  "name": "example",
		  "datasetSizeInGb": 1,
		  "ramPercentage": 40,
		  "supportOSSClusterApi": false,
		  "respVersion": "resp3",
		  "useExternalEndpointForOSSClusterApi": false,
		  "dataEvictionPolicy": "allkeys-lru",
		  "replication": true,
		  "throughputMeasurement": {
		    "by": "operations-per-second",
		    "value": 1000
		  },
		  "regexRules": [".*"],
		  "dataPersistence": "none",
		  "replicaOf": [
		    "another"
		  ],
		  "periodicBackupPath": "s3://bucket-name",
		  "sourceIp": [
		    "10.0.0.1"
		  ],
		  "clientSslCertificate": "something",
		  "clientTlsCertificates": ["something", "new"],
		  "enableTls": false,
		  "password": "fooBar",
		  "alerts": [
		    {
		      "name": "dataset-size",
		      "value": 80
		    }
		  ],
		  "enableDefaultUser": false,
		  "queryPerformanceFactor": "2x",
		  "autoMinorVersionUpgrade": true
		}`,
		"task",
		"databaseUpdateRequest",
	)

	s := httptest.NewServer(testServer("key", "secret", flow...))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Database.Update(context.TODO(), 42, 18, databases.UpdateDatabase{
		DryRun:                              redis.Bool(false),
		Name:                                redis.String("example"),
		DatasetSizeInGB:                     redis.Float64(1),
		RamPercentage:                       redis.Int(40),
		SupportOSSClusterAPI:                redis.Bool(false),
		RespVersion:                         redis.String("resp3"),
		UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
		DataPersistence:                     redis.String("none"),
		DataEvictionPolicy:                  redis.String("allkeys-lru"),
		Replication:                         redis.Bool(true),
		ThroughputMeasurement: &databases.UpdateThroughputMeasurement{
			By:    redis.String("operations-per-second"),
			Value: redis.Int(1000),
		},
		RegexRules:            redis.StringSlice(".*"),
		ReplicaOf:             redis.StringSlice("another"),
		PeriodicBackupPath:    redis.String("s3://bucket-name"),
		SourceIP:              redis.StringSlice("10.0.0.1"),
		ClientSSLCertificate:  redis.String("something"),
		ClientTLSCertificates: &[]*string{redis.String("something"), redis.String("new")},
		EnableTls:             redis.Bool(false),
		Password:              redis.String("fooBar"),
		Alerts: &[]*databases.Alert{
			{
				Name:  redis.String("dataset-size"),
				Value: redis.Int(80),
			},
		},
		EnableDefaultUser:       redis.Bool(false),
		QueryPerformanceFactor:  redis.String("2x"),
		AutoMinorVersionUpgrade: redis.Bool(true),
	})
	require.NoError(t, err)
}

func TestDatabase_Delete(t *testing.T) {
	flow := taskFlow(
		t,
		http.MethodDelete,
		"/subscriptions/42/databases/4291",
		"",
		"task",
		"databaseDeleteRequest",
	)

	s := httptest.NewServer(testServer("key", "secret", flow...))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Database.Delete(context.TODO(), 42, 4291)
	require.NoError(t, err)
}

func TestDatabase_Backup(t *testing.T) {
	flow := taskFlow(
		t,
		http.MethodPost,
		"/subscriptions/42/databases/18/backup",
		"",
		"task-uuid",
		"databaseBackupRequest",
	)

	s := httptest.NewServer(testServer("key", "secret", flow...))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Database.Backup(context.TODO(), 42, 18)
	require.NoError(t, err)
}

func TestDatabase_Import(t *testing.T) {
	flow := taskFlow(
		t,
		http.MethodPost,
		"/subscriptions/42/databases/81/import",
		`{
		  "sourceType": "magic",
		  "importFromUri": ["tinkerbell"]
		}`,
		"task-uuid",
		"databaseImportRequest",
	)

	s := httptest.NewServer(testServer("key", "secret", flow...))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Database.Import(context.TODO(), 42, 81, databases.Import{
		SourceType:    redis.String("magic"),
		ImportFromURI: redis.StringSlice("tinkerbell"),
	})
	require.NoError(t, err)
}

func TestDatabase_Certificate(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", getRequest(t, "/subscriptions/42/databases/18/certificate",
		`{ "publicCertificatePEMString": "public-cert" }`)))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	certificate, err := subject.Database.GetCertificate(context.TODO(), 42, 18)
	require.NoError(t, err)

	assert.Equal(t, &databases.DatabaseCertificate{
		PublicCertificatePEMString: "public-cert",
	}, certificate)

}

func TestDatabase_UpgradeRedisVersion(t *testing.T) {
	flow := taskFlow(
		t,
		http.MethodPost,
		"/subscriptions/42/databases/18/upgrade",
		`{ "targetRedisVersion": "7.2" }`,
		"upgrade-task-id",
		"databaseUpgradeRequest",
	)

	s := httptest.NewServer(testServer("key", "secret", flow...))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Database.UpgradeRedisVersion(
		context.TODO(),
		42,
		18,
		databases.UpgradeRedisVersion{
			TargetRedisVersion: redis.String("7.2"),
		},
	)
	require.NoError(t, err)
}
