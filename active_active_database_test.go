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

func TestAADatabase_Create(t *testing.T) {
	expected := 1466
	s := httptest.NewServer(testServer("key", "secret", postRequest(t, "/subscriptions/111478/databases", `{
  "dryRun": false,
  "name": "active-active-example",
  "protocol": "redis",
  "memoryLimitInGb": 1,
  "datasetSizeInGb": 1,
  "supportOSSClusterApi": false,
  "respVersion": "resp3",
  "useExternalEndpointForOSSClusterApi": false,
  "dataEvictionPolicy": "noeviction",
  "dataPersistence": "none",
  "sourceIp": [
    "0.0.0.0/0"
  ],
  "password": "test-password",
  "alerts": [
    {
      "name": "dataset-size",
      "value": 80
    }
  ],
  "modules": [
    {
      "name": "RedisJSON"
    }
  ],
  "localThroughputMeasurement": [
    {
      "region": "us-east-1",
      "writeOperationsPerSecond": 1000,
      "readOperationsPerSecond": 1000
    },
    {
      "region": "us-east-2",
      "writeOperationsPerSecond": 1000,
      "readOperationsPerSecond": 1000
    }
  ],
  "port": 12000,
  "queryPerformanceFactor": "Standard",
  "redisVersion": "7.2",
  "autoMinorVersionUpgrade": true
}`, `{
  "taskId": "task",
  "commandType": "databaseCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2024-05-08T08:10:02Z",
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
  "timestamp": "2024-05-08T08:22:34Z",
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

	actual, err := subject.Database.ActiveActiveCreate(context.TODO(), 111478, databases.CreateActiveActiveDatabase{
		DryRun:                              redis.Bool(false),
		Name:                                redis.String("active-active-example"),
		Protocol:                            redis.String("redis"),
		MemoryLimitInGB:                     redis.Float64(1),
		DatasetSizeInGB:                     redis.Float64(1),
		SupportOSSClusterAPI:                redis.Bool(false),
		RespVersion:                         redis.String("resp3"),
		UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
		DataEvictionPolicy:                  redis.String("noeviction"),
		GlobalDataPersistence:               redis.String("none"),
		GlobalSourceIP:                      redis.StringSlice("0.0.0.0/0"),
		GlobalPassword:                      redis.String("test-password"),
		GlobalAlerts: []*databases.Alert{
			{
				Name:  redis.String("dataset-size"),
				Value: redis.Int(80),
			},
		},
		GlobalModules: []*databases.Module{
			{
				Name: redis.String("RedisJSON"),
			},
		},
		LocalThroughputMeasurement: []*databases.LocalThroughput{
			{
				Region:                   redis.String("us-east-1"),
				WriteOperationsPerSecond: redis.Int(1000),
				ReadOperationsPerSecond:  redis.Int(1000),
			},
			{
				Region:                   redis.String("us-east-2"),
				WriteOperationsPerSecond: redis.Int(1000),
				ReadOperationsPerSecond:  redis.Int(1000),
			},
		},
		PortNumber:              redis.Int(12000),
		QueryPerformanceFactor:  redis.String("Standard"),
		RedisVersion:            redis.String("7.2"),
		AutoMinorVersionUpgrade: redis.Bool(true),
	})
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestAADatabase_Update(t *testing.T) {
	flow := taskFlow(
		t,
		http.MethodPut,
		"/subscriptions/111478/databases/1466/regions",
		`{
  "dryRun": false,
  "memoryLimitInGb": 2,
  "datasetSizeInGb": 2,
  "supportOSSClusterApi": false,
  "useExternalEndpointForOSSClusterApi": false,
  "clientSslCertificate": "cert-content",
  "clientTlsCertificates": ["cert1", "cert2"],
  "enableTls": true,
  "globalDataPersistence": "aof-every-1-second",
  "globalPassword": "new-password",
  "globalEnableDefaultUser": true,
  "globalSourceIp": [
    "192.168.1.0/24"
  ],
  "globalAlerts": [
    {
      "name": "throughput-higher-than",
      "value": 90
    }
  ],
  "regions": [
    {
      "region": "us-east-1",
      "remoteBackup": {
        "active": true,
        "interval": "every-12-hours",
        "timeUTC": "10:00",
        "storageType": "aws-s3",
        "storagePath": "s3://bucket/path"
      },
      "localThroughputMeasurement": {
        "writeOperationsPerSecond": 2000,
        "readOperationsPerSecond": 2000
      },
      "dataPersistence": "aof-every-1-second",
      "password": "region-password",
      "sourceIp": [
        "10.0.0.0/8"
      ],
      "enableDefaultUser": false,
      "alerts": [
        {
          "name": "dataset-size",
          "value": 85
        }
      ]
    }
  ],
  "dataEvictionPolicy": "allkeys-lru",
  "queryPerformanceFactor": "6x",
  "autoMinorVersionUpgrade": true
}`,
		"update-task",
		"databaseUpdateRequest",
	)

	s := httptest.NewServer(testServer("key", "secret", flow...))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Database.ActiveActiveUpdate(context.TODO(), 111478, 1466, databases.UpdateActiveActiveDatabase{
		DryRun:                              redis.Bool(false),
		MemoryLimitInGB:                     redis.Float64(2),
		DatasetSizeInGB:                     redis.Float64(2),
		SupportOSSClusterAPI:                redis.Bool(false),
		UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
		ClientSSLCertificate:                redis.String("cert-content"),
		ClientTLSCertificates:               &[]*string{redis.String("cert1"), redis.String("cert2")},
		EnableTls:                           redis.Bool(true),
		GlobalDataPersistence:               redis.String("aof-every-1-second"),
		GlobalPassword:                      redis.String("new-password"),
		GlobalEnableDefaultUser:             redis.Bool(true),
		GlobalSourceIP:                      redis.StringSlice("192.168.1.0/24"),
		GlobalAlerts: &[]*databases.Alert{
			{
				Name:  redis.String("throughput-higher-than"),
				Value: redis.Int(90),
			},
		},
		Regions: []*databases.LocalRegionProperties{
			{
				Region: redis.String("us-east-1"),
				RemoteBackup: &databases.DatabaseBackupConfig{
					Active:      redis.Bool(true),
					Interval:    redis.String("every-12-hours"),
					TimeUTC:     redis.String("10:00"),
					StorageType: redis.String("aws-s3"),
					StoragePath: redis.String("s3://bucket/path"),
				},
				LocalThroughputMeasurement: &databases.LocalThroughput{
					WriteOperationsPerSecond: redis.Int(2000),
					ReadOperationsPerSecond:  redis.Int(2000),
				},
				DataPersistence:   redis.String("aof-every-1-second"),
				Password:          redis.String("region-password"),
				SourceIP:          redis.StringSlice("10.0.0.0/8"),
				EnableDefaultUser: redis.Bool(false),
				Alerts: &[]*databases.Alert{
					{
						Name:  redis.String("dataset-size"),
						Value: redis.Int(85),
					},
				},
			},
		},
		DataEvictionPolicy:      redis.String("allkeys-lru"),
		QueryPerformanceFactor:  redis.String("6x"),
		AutoMinorVersionUpgrade: redis.Bool(true),
	})
	require.NoError(t, err)
}

func TestAADatabase_List(t *testing.T) {
	body := `{
					"accountId": 69369,
					"subscription": [
						{
							"subscriptionId": 111478,
							"numberOfDatabases": 1,
							"databases": [
								{
									"databaseId": 1466,
									"name": "creation-plan-db-1",
									"redisVersion": "7.2",
									"protocol": "redis",
									"status": "active",
									"memoryStorage": "ram",
									"activeActiveRedis": true,
									"activatedOn": "2024-05-08T08:10:02Z",
									"lastModified": "2024-05-08T08:22:34Z",
									"supportOSSClusterApi": false,
									"useExternalEndpointForOSSClusterApi": false,
									"replication": true,
									"dataEvictionPolicy": "noeviction",
									"autoMinorVersionUpgrade": true,
									"modules": [],
									"globalDataPersistence": "none",
									"globalSourceIp": ["0.0.0.0/0"],
									"globalAlerts": [
										{
											"name": "dataset-size",
											"value": 80
										}
									],
									"globalModules": [],
									"globalEnableDefaultUser": true,
									"crdbDatabases": [
										{
											"provider": "AWS",
											"region": "us-east-1",
											"redisVersionCompliance": "6.2.10",
											"respVersion": "resp2",
											"publicEndpoint": "redis-14383.mc940-1.us-east-1-mz.ec2.qa-cloud.rlrcp.com:14383",
											"privateEndpoint": "redis-14383.internal.mc940-1.us-east-1-mz.ec2.qa-cloud.rlrcp.com:14383",
											"memoryLimitInGb": 1,
											"datasetSizeInGb": 1,
											"memoryUsedInMb": 29.9949,
											"readOperationsPerSecond": 1000,
											"writeOperationsPerSecond": 1000,
											"dataPersistence": "none",
											"queryPerformanceFactor": "Standard",
											"alerts": [
												{
													"id": "51054122-2",
													"name": "dataset-size",
													"value": 80,
													"defaultValue": 80
												}
											],
											"security": {
												"enableDefaultUser": true,
												"sslClientAuthentication": false,
												"tlsClientAuthentication": false,
												"enableTls": false,
												"sourceIps": [
													"0.0.0.0/0"
												]
											},
											"backup": {
												"enableRemoteBackup": false,
												"interval": "",
												"timeUTC": ""
											},
											"links": []
										},
										{
											"provider": "AWS",
											"region": "us-east-2",
											"redisVersionCompliance": "6.2.10",
											"respVersion": "resp2",
											"publicEndpoint": "redis-14383.mc940-0.us-east-2-mz.ec2.qa-cloud.rlrcp.com:14383",
											"privateEndpoint": "redis-14383.internal.mc940-0.us-east-2-mz.ec2.qa-cloud.rlrcp.com:14383",
											"memoryLimitInGb": 1,
											"datasetSizeInGb": 1,
											"memoryUsedInMb": 29.9788,
											"readOperationsPerSecond": 1000,
											"writeOperationsPerSecond": 1000,
											"dataPersistence": "none",
											"queryPerformanceFactor": "Standard",
											"alerts": [
												{
													"id": "51054121-2",
													"name": "dataset-size",
													"value": 80,
													"defaultValue": 80
												}
											],
											"security": {
												"enableDefaultUser": true,
												"sslClientAuthentication": false,
												"tlsClientAuthentication": false,
												"enableTls": false,
												"sourceIps": [
													"0.0.0.0/0"
												]
											},
											"backup": {
												"enableRemoteBackup": false,
												"interval": "",
												"timeUTC": ""
											},
											"links": []
										}
									],
									"links": []
								}
							],
							"links": []
						}
					],
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/subscriptions/111478/databases?offset=0&limit=100"
						}
					]
				}`

	query := map[string][]string{"limit": {"100"}, "offset": {"0"}}

	s := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			getRequestWithQuery(
				t,
				"/subscriptions/111478/databases",
				query,
				body,
			),
			getRequestWithQueryAndStatus(
				t,
				"/subscriptions/111478/databases",
				map[string][]string{"limit": {"100"}, "offset": {"100"}},
				404,
				"",
			),
		),
	)

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	list := subject.Database.ListActiveActive(context.TODO(), 111478)

	var actual []*databases.ActiveActiveDatabase
	for list.Next() {
		actual = append(actual, list.Value())
	}
	require.NoError(t, list.Err())

	assert.Equal(t, []*databases.ActiveActiveDatabase{
		{
			ID:                                  redis.Int(1466),
			Name:                                redis.String("creation-plan-db-1"),
			Protocol:                            redis.String("redis"),
			RedisVersion:                        redis.String("7.2"),
			Status:                              redis.String("active"),
			MemoryStorage:                       redis.String("ram"),
			ActiveActiveRedis:                   redis.Bool(true),
			ActivatedOn:                         redis.Time(time.Date(2024, 5, 8, 8, 10, 02, 0, time.UTC)),
			LastModified:                        redis.Time(time.Date(2024, 5, 8, 8, 22, 34, 0, time.UTC)),
			SupportOSSClusterAPI:                redis.Bool(false),
			UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
			Replication:                         redis.Bool(true),
			DataEvictionPolicy:                  redis.String("noeviction"),
			AutoMinorVersionUpgrade:             redis.Bool(true),
			Modules:                             []*databases.Module{},
			GlobalDataPersistence:               redis.String("none"),
			GlobalSourceIP:                      redis.StringSlice("0.0.0.0/0"),
			GlobalAlerts: []*databases.Alert{
				{
					Name:  redis.String("dataset-size"),
					Value: redis.Int(80),
				},
			},
			GlobalModules:           []*databases.Module{},
			GlobalEnableDefaultUser: redis.Bool(true),
			CrdbDatabases: []*databases.CrdbDatabase{
				{
					Provider:                 redis.String("AWS"),
					Region:                   redis.String("us-east-1"),
					RedisVersionCompliance:   redis.String("6.2.10"),
					PublicEndpoint:           redis.String("redis-14383.mc940-1.us-east-1-mz.ec2.qa-cloud.rlrcp.com:14383"),
					PrivateEndpoint:          redis.String("redis-14383.internal.mc940-1.us-east-1-mz.ec2.qa-cloud.rlrcp.com:14383"),
					MemoryLimitInGB:          redis.Float64(1),
					DatasetSizeInGB:          redis.Float64(1),
					MemoryUsedInMB:           redis.Float64(29.9949),
					ReadOperationsPerSecond:  redis.Int(1000),
					WriteOperationsPerSecond: redis.Int(1000),
					DataPersistence:          redis.String("none"),
					QueryPerformanceFactor:   redis.String("Standard"),
					Alerts: []*databases.Alert{
						{
							Name:  redis.String("dataset-size"),
							Value: redis.Int(80),
						},
					},
					Security: &databases.Security{
						EnableDefaultUser:       redis.Bool(true),
						SSLClientAuthentication: redis.Bool(false),
						TLSClientAuthentication: redis.Bool(false),
						SourceIPs:               redis.StringSlice("0.0.0.0/0"),
						EnableTls:               redis.Bool(false),
					},
					Backup: &databases.Backup{
						Enabled:  redis.Bool(false),
						TimeUTC:  redis.String(""),
						Interval: redis.String(""),
					},
				},
				{
					Provider:                 redis.String("AWS"),
					Region:                   redis.String("us-east-2"),
					RedisVersionCompliance:   redis.String("6.2.10"),
					PublicEndpoint:           redis.String("redis-14383.mc940-0.us-east-2-mz.ec2.qa-cloud.rlrcp.com:14383"),
					PrivateEndpoint:          redis.String("redis-14383.internal.mc940-0.us-east-2-mz.ec2.qa-cloud.rlrcp.com:14383"),
					MemoryLimitInGB:          redis.Float64(1),
					DatasetSizeInGB:          redis.Float64(1),
					MemoryUsedInMB:           redis.Float64(29.9788),
					ReadOperationsPerSecond:  redis.Int(1000),
					WriteOperationsPerSecond: redis.Int(1000),
					DataPersistence:          redis.String("none"),
					QueryPerformanceFactor:   redis.String("Standard"),
					Alerts: []*databases.Alert{
						{
							Name:  redis.String("dataset-size"),
							Value: redis.Int(80),
						},
					},
					Security: &databases.Security{
						EnableDefaultUser:       redis.Bool(true),
						SSLClientAuthentication: redis.Bool(false),
						TLSClientAuthentication: redis.Bool(false),
						SourceIPs:               redis.StringSlice("0.0.0.0/0"),
						EnableTls:               redis.Bool(false),
					},
					Backup: &databases.Backup{
						Enabled:  redis.Bool(false),
						TimeUTC:  redis.String(""),
						Interval: redis.String(""),
					},
				},
			},
		},
	}, actual)

}
