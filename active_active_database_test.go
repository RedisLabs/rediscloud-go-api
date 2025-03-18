package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
									"modules": [],
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
			Status:                              redis.String("active"),
			MemoryStorage:                       redis.String("ram"),
			ActiveActiveRedis:                   redis.Bool(true),
			ActivatedOn:                         redis.Time(time.Date(2024, 5, 8, 8, 10, 02, 0, time.UTC)),
			LastModified:                        redis.Time(time.Date(2024, 5, 8, 8, 22, 34, 0, time.UTC)),
			SupportOSSClusterAPI:                redis.Bool(false),
			UseExternalEndpointForOSSClusterAPI: redis.Bool(false),
			Replication:                         redis.Bool(true),
			DataEvictionPolicy:                  redis.String("noeviction"),
			Modules:                             []*databases.Module{},
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
