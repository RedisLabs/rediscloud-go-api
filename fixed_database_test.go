package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	fixedDatabases "github.com/RedisLabs/rediscloud-go-api/service/fixed/databases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedDatabase_Create(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			postRequest(
				t,
				"/fixed/subscriptions/111728/databases",
				`{
					"name": "my-test-fixed-database",
					"protocol": "memcached",
					"respVersion": "resp2",
					"dataPersistence": "none",
					"dataEvictionPolicy": "noeviction",
					"replication": false,
					"alerts": []
				}`,
				`{
					"taskId": "784299af-17ea-4ed6-b08f-dd643238c8dd",
					"commandType": "fixedDatabaseCreateRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-10T14:14:14.736763484Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/784299af-17ea-4ed6-b08f-dd643238c8dd"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/784299af-17ea-4ed6-b08f-dd643238c8dd",
				`{
					"taskId": "784299af-17ea-4ed6-b08f-dd643238c8dd",
					"commandType": "fixedDatabaseCreateRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-10T14:14:34.153537279Z",
					"response": {
						"resourceId": 51055029,
						"additionalResourceId": 111728
					},
					"links": [
						{
							"rel": "resource",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111728/databases/51055029"
						},
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/784299af-17ea-4ed6-b08f-dd643238c8dd"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedDatabases.Create(
		context.TODO(),
		111728,
		fixedDatabases.CreateFixedDatabase{
			Name:               redis.String("my-test-fixed-database"),
			Protocol:           redis.String("memcached"),
			RespVersion:        redis.String("resp2"),
			DataPersistence:    redis.String("none"),
			DataEvictionPolicy: redis.String("noeviction"),
			Replication:        redis.Bool(false),
			Alerts:             &[]*databases.Alert{},
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 51055029, actual)
}

func TestFixedDatabase_Create_with_RedisVersion(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			postRequest(
				t,
				"/fixed/subscriptions/111728/databases",
				`{
					"name": "my-redis-essentials-db",
					"protocol": "redis",
					"redisVersion": "7.4",
					"dataPersistence": "none",
					"dataEvictionPolicy": "noeviction",
					"replication": false,
					"alerts": []
				}`,
				`{
					"taskId": "784299af-17ea-4ed6-b08f-dd643238c8dd",
					"commandType": "fixedDatabaseCreateRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-10T14:14:14.736763484Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/784299af-17ea-4ed6-b08f-dd643238c8dd"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/784299af-17ea-4ed6-b08f-dd643238c8dd",
				`{
					"taskId": "784299af-17ea-4ed6-b08f-dd643238c8dd",
					"commandType": "fixedDatabaseCreateRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-10T14:14:34.153537279Z",
					"response": {
						"resourceId": 51055030,
						"additionalResourceId": 111728
					},
					"links": [
						{
							"rel": "resource",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111728/databases/51055030"
						},
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/784299af-17ea-4ed6-b08f-dd643238c8dd"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedDatabases.Create(
		context.TODO(),
		111728,
		fixedDatabases.CreateFixedDatabase{
			Name:               redis.String("my-redis-essentials-db"),
			Protocol:           redis.String("redis"),
			RedisVersion:       redis.String("7.4"),
			DataPersistence:    redis.String("none"),
			DataEvictionPolicy: redis.String("noeviction"),
			Replication:        redis.Bool(false),
			Alerts:             &[]*databases.Alert{},
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 51055030, actual)
}

func TestFixedDatabase_List(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			getRequestWithQuery(
				t,
				"/fixed/subscriptions/111930/databases",
				map[string][]string{
					"limit": {
						"100",
					},
					"offset": {
						"0",
					},
				},
				`{
					"accountId": 69369,
					"subscription": {
						"subscriptionId": 111930,
						"numberOfDatabases": 1,
						"databases": [
							{
								"databaseId": 51055698,
								"name": "my-second-test-fixed-database",
								"protocol": "memcached",
								"provider": "AWS",
								"region": "us-west-1",
								"respVersion": "resp2",
								"status": "draft",
								"planMemoryLimit": 1,
								"memoryLimitMeasurementUnit": "GB",
								"memoryUsedInMb": 7,
								"memoryStorage": "ram",
								"supportOSSClusterApi": false,
								"useExternalEndpointForOSSClusterApi": false,
								"dataPersistence": "none",
								"replication": false,
								"dataEvictionPolicy": "noeviction",
								"activatedOn": "2024-05-14T09:27:48Z",
								"replicaOf": null,
								"replica": null,
								"clustering": {
									"enabled": true,
									"regexRules": [
										{
											"ordinal": 0,
											"pattern": ".*\\{(?<tag>.*)\\}.*"
										},
										{
											"ordinal": 1,
											"pattern": "(?<tag>.*)"
										}
									],
									"hashingPolicy": "standard"
								},
								"security": {
									"defaultUserEnabled": true,
									"sslClientAuthentication": false,
									"tlsClientAuthentication": false,
									"enableTls": false,
									"sourceIps": [
										"0.0.0.0/0"
									]
								},
								"modules": [],
								"alerts": [],
								"backup": {
									"remoteBackupEnabled": false,
									"status": "idle"
								},
								"links": []
							}
						],
						"links": []
					},
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111930/databases?offset=0&limit=100"
						}
					]
				}`,
			),
			getRequestWithQuery(
				t,
				"/fixed/subscriptions/111930/databases",
				map[string][]string{
					"limit": {
						"100",
					},
					"offset": {
						"100",
					},
				},
				`{
					"accountId" : 69369,
					"subscription" : {
						"subscriptionId" : 111930,
						"numberOfDatabases" : 0,
						"databases" : [ ],
						"links" : [ ]
					},
					"links" : [ {
						"rel" : "self",
						"type" : "GET",
						"href" : "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/112330/databases?limit=100&offset=100"
					} ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	list := subject.FixedDatabases.List(context.TODO(), 111930)

	var actual []*fixedDatabases.FixedDatabase
	for list.Next() {
		actual = append(actual, list.Value())
	}
	require.NoError(t, list.Err())

	assert.Equal(t, []*fixedDatabases.FixedDatabase{
		{
			DatabaseId:                          redis.Int(51055698),
			Name:                                redis.String("my-second-test-fixed-database"),
			Protocol:                            redis.String("memcached"),
			Provider:                            redis.String("AWS"),
			Region:                              redis.String("us-west-1"),
			RespVersion:                         redis.String("resp2"),
			Status:                              redis.String("draft"),
			PlanMemoryLimit:                     redis.Float64(1),
			MemoryLimitMeasurementUnit:          redis.String("GB"),
			MemoryUsedInMb:                      redis.Float64(7),
			MemoryStorage:                       redis.String("ram"),
			SupportOSSClusterAPI:                redis.Bool(false),
			UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
			DataPersistence:                     redis.String("none"),
			Replication:                         redis.Bool(false),
			DataEvictionPolicy:                  redis.String("noeviction"),
			ActivatedOn:                         redis.Time(time.Date(2024, 5, 14, 9, 27, 48, 0, time.UTC)),
			Clustering: &fixedDatabases.Clustering{
				Enabled: redis.Bool(true),
				RegexRules: []*databases.RegexRule{
					{
						Ordinal: 0,
						Pattern: ".*\\{(?<tag>.*)\\}.*",
					},
					{
						Ordinal: 1,
						Pattern: "(?<tag>.*)",
					},
				},
				HashingPolicy: redis.String("standard"),
			},
			Security: &fixedDatabases.Security{
				EnableDefaultUser:       redis.Bool(true),
				SSLClientAuthentication: redis.Bool(false),
				TLSClientAuthentication: redis.Bool(false),
				EnableTls:               redis.Bool(false),
				SourceIPs:               redis.StringSlice("0.0.0.0/0"),
			},
			Modules: &[]*databases.Module{},
			Alerts:  &[]*databases.Alert{},
			Backup: &fixedDatabases.Backup{
				Enabled: redis.Bool(false),
				Status:  redis.String("idle"),
			},
		},
	}, actual)

}

func TestFixedDatabase_Get(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/111728/databases/51055029",
				`{
					"databaseId": 51055029,
					"name": "my-test-fixed-database",
					"protocol": "memcached",
					"provider": "AWS",
					"region": "us-west-1",
					"respVersion": "resp2",
					"status": "draft",
					"planMemoryLimit": 1,
					"memoryLimitMeasurementUnit": "GB",
					"memoryUsedInMb": 7,
					"memoryStorage": "ram",
					"supportOSSClusterApi": false,
					"useExternalEndpointForOSSClusterApi": false,
					"dataPersistence": "none",
					"replication": false,
					"dataEvictionPolicy": "noeviction",
					"activatedOn": "2024-05-10T14:14:33Z",
					"replicaOf": null,
					"replica": null,
					"clustering": {
					"enabled": true,
						"regexRules": [
							{
								"ordinal": 0,
								"pattern": ".*\\{(?<tag>.*)\\}.*"
							},
							{
								"ordinal": 1,
								"pattern": "(?<tag>.*)"
							}
						],
						"hashingPolicy": "standard"
					},
					"security": {
						"defaultUserEnabled": true,
						"password": "s1j5CSgTliqM0EHzMU68GbflSVGuCIcF",
						"sslClientAuthentication": false,
						"tlsClientAuthentication": false,
						"enableTls": false,
						"sourceIps": [
							"0.0.0.0/0"
						]
					},
					"modules": [],
					"alerts": [],
					"backup": {
						"remoteBackupEnabled": false,
						"status": "idle"
					},
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111728/databases/51055029"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedDatabases.Get(context.TODO(), 111728, 51055029)
	require.NoError(t, err)

	assert.Equal(t, &fixedDatabases.FixedDatabase{
		DatabaseId:                          redis.Int(51055029),
		Name:                                redis.String("my-test-fixed-database"),
		Protocol:                            redis.String("memcached"),
		Provider:                            redis.String("AWS"),
		Region:                              redis.String("us-west-1"),
		RespVersion:                         redis.String("resp2"),
		Status:                              redis.String("draft"),
		PlanMemoryLimit:                     redis.Float64(1),
		MemoryLimitMeasurementUnit:          redis.String("GB"),
		MemoryUsedInMb:                      redis.Float64(7),
		MemoryStorage:                       redis.String("ram"),
		SupportOSSClusterAPI:                redis.Bool(false),
		UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
		DataPersistence:                     redis.String("none"),
		Replication:                         redis.Bool(false),
		DataEvictionPolicy:                  redis.String("noeviction"),
		ActivatedOn:                         redis.Time(time.Date(2024, 5, 10, 14, 14, 33, 0, time.UTC)),
		Clustering: &fixedDatabases.Clustering{
			Enabled: redis.Bool(true),
			RegexRules: []*databases.RegexRule{
				{
					Ordinal: 0,
					Pattern: ".*\\{(?<tag>.*)\\}.*",
				},
				{
					Ordinal: 1,
					Pattern: "(?<tag>.*)",
				},
			},
			HashingPolicy: redis.String("standard"),
		},
		Security: &fixedDatabases.Security{
			EnableDefaultUser:       redis.Bool(true),
			Password:                redis.String("s1j5CSgTliqM0EHzMU68GbflSVGuCIcF"),
			SSLClientAuthentication: redis.Bool(false),
			TLSClientAuthentication: redis.Bool(false),
			EnableTls:               redis.Bool(false),
			SourceIPs:               redis.StringSlice("0.0.0.0/0"),
		},
		Modules: &[]*databases.Module{},
		Alerts:  &[]*databases.Alert{},
		Backup: &fixedDatabases.Backup{
			Enabled: redis.Bool(false),
			Status:  redis.String("idle"),
		},
	}, actual)

}

func TestFixedDatabase_Update(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			putRequest(
				t,
				"/fixed/subscriptions/112119/databases/51056892",
				`{
					"name": "my-test-fixed-database",
					"respVersion": "resp2",
					"dataPersistence": "none",
					"dataEvictionPolicy": "volatile-lru",
					"replication": false,
					"enableDefaultUser": true,
					"alerts": [
						{
							"name": "datasets-size",
							"value": 80
						}
					]
				}`,
				`{
					"taskId": "4a9dc555-7c7d-4657-a8c9-5d868d7acdd2",
					"commandType": "fixedDatabaseUpdateRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-15T14:47:08.975576016Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/4a9dc555-7c7d-4657-a8c9-5d868d7acdd2"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/4a9dc555-7c7d-4657-a8c9-5d868d7acdd2",
				`{
					"taskId": "4a9dc555-7c7d-4657-a8c9-5d868d7acdd2",
					"commandType": "fixedDatabaseUpdateRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-15T14:47:21.687621578Z",
					"response": {
						"resourceId": 51056892,
						"additionalResourceId": 112119
					},
					"links": [
						{
							"rel": "resource",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/112119/databases/51056892"
						},
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/4a9dc555-7c7d-4657-a8c9-5d868d7acdd2"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.FixedDatabases.Update(
		context.TODO(),
		112119,
		51056892,
		fixedDatabases.UpdateFixedDatabase{
			Name:               redis.String("my-test-fixed-database"),
			RespVersion:        redis.String("resp2"),
			DataPersistence:    redis.String("none"),
			DataEvictionPolicy: redis.String("volatile-lru"),
			Replication:        redis.Bool(false),
			EnableDefaultUser:  redis.Bool(true),
			Alerts: &[]*databases.Alert{
				{
					Name:  redis.String("datasets-size"),
					Value: redis.Int(80),
				},
			},
		},
	)

	require.NoError(t, err)
}

func TestFixedDatabase_Delete(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			deleteRequest(
				t,
				"/fixed/subscriptions/112119/databases/51056892",
				`{
					"taskId": "3e26ac3a-231c-486d-b4cf-8519f520a5f4",
					"commandType": "fixedDatabaseDeleteRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-15T14:55:04.008723915Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/3e26ac3a-231c-486d-b4cf-8519f520a5f4"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/3e26ac3a-231c-486d-b4cf-8519f520a5f4",
				`{
					"taskId": "3e26ac3a-231c-486d-b4cf-8519f520a5f4",
					"commandType": "fixedDatabaseDeleteRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-15T14:55:06.538386979Z",
					"response": {
						"resourceId": 51056892,
						"additionalResourceId": 112119
					},
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/3e26ac3a-231c-486d-b4cf-8519f520a5f4"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.FixedDatabases.Delete(context.TODO(), 112119, 51056892)
	require.NoError(t, err)

}

func TestFixedDatabase_UpgradeRedisVersion(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			postRequest(
				t,
				"/fixed/subscriptions/112119/databases/51056892/upgrade",
				`{ "targetRedisVersion": "7.4" }`,
				`{
					"taskId": "upgrade-task-uuid",
					"commandType": "fixedDatabaseUpgradeRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-15T14:55:04.008723915Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/upgrade-task-uuid"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/upgrade-task-uuid",
				`{
					"taskId": "upgrade-task-uuid",
					"commandType": "fixedDatabaseUpgradeRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-15T14:55:06.538386979Z",
					"response": {
						"resourceId": 51056892,
						"additionalResourceId": 112119
					},
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/upgrade-task-uuid"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.FixedDatabases.UpgradeRedisVersion(
		context.TODO(),
		112119,
		51056892,
		fixedDatabases.UpgradeFixedDatabaseRedisVersion{
			TargetRedisVersion: redis.String("7.4"),
		},
	)

	require.NoError(t, err)
}
