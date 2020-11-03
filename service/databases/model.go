package databases

import (
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type taskResponse struct {
	ID *string `json:"taskId,omitempty"`
}

func (o taskResponse) String() string {
	return internal.ToString(o)
}

type Database struct {
	ID                     *int        `json:"databaseId,omitempty"`
	Name                   *string     `json:"name,omitempty"`
	Protocol               *string     `json:"protocol,omitempty"`
	Provider               *string     `json:"provider,omitempty"`
	Region                 *string     `json:"region,omitempty"`
	Status                 *string     `json:"status,omitempty"`
	MemoryLimitInGb        *float64    `json:"memoryLimitInGb,omitempty"`
	MemoryUsedInMb         *float64    `json:"memoryUsedInMb,omitempty"`
	SupportOSSClusterApi   *bool       `json:"supportOSSClusterApi,omitempty"`
	DataPersistence        *string     `json:"dataPersistence,omitempty"`
	Replication            *bool       `json:"replication,omitempty"`
	DataEvictionPolicy     *string     `json:"dataEvictionPolicy,omitempty"`
	ThroughputMeasurement  *Throughput `json:"throughputMeasurement,omitempty"`
	ReplicaOf              []*string   `json:"replicaOf,omitempty"`
	Clustering             *Clustering `json:"clustering,omitempty"`
	Security               *Security   `json:"security,omitempty"`
	Modules                []*Module   `json:"modules,omitempty"`
	Alerts                 []*Alert    `json:"alerts,omitempty"`
	ActivatedOn            *time.Time  `json:"activatedOn,omitempty"`
	LastModified           *time.Time  `json:"lastModified,omitempty"`
	MemoryStorage          *string     `json:"memoryStorage,omitempty"`
	PrivateEndpoint        *string     `json:"privateEndpoint,omitempty"`
	PublicEndpoint         *string     `json:"publicEndpoint,omitempty"`
	RedisVersionCompliance *string     `json:"redisVersionCompliance,omitempty"`
}

func (o Database) String() string {
	return internal.ToString(o)
}

type Clustering struct {
	NumberOfShards *int `json:"numberOfShards,omitempty"`
	// TODO RegexRules interface{} `json:"regexRules,omitempty"`
	// TODO HashingPolicy interface{} `json:"hashingPolicy,omitempty"`
}

func (o Clustering) String() string {
	return internal.ToString(o)
}

type Security struct {
	SslClientAuthentication *bool     `json:"sslClientAuthentication,omitempty"`
	SourceIps               []*string `json:"sourceIps,omitempty"`
	Password                *string   `json:"password,omitempty"`
}

func (o Security) String() string {
	return internal.ToString(o)
}

type Module struct {
	Name       *string            `json:"name,omitempty"`
	Parameters map[string]*string `json:"parameters,omitempty"`
}

func (o Module) String() string {
	return internal.ToString(o)
}

type Throughput struct {
	By    *string `json:"by,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o Throughput) String() string {
	return internal.ToString(o)
}

type Alert struct {
	Name  *string `json:"name,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o Alert) String() string {
	return internal.ToString(o)
}

type listDatabaseResponse struct {
	Subscription []*listDbSubscription `json:"subscription,omitempty"`
}

func (o listDatabaseResponse) String() string {
	return internal.ToString(o)
}

type listDbSubscription struct {
	ID        *int        `json:"subscriptionId,omitempty"`
	Databases []*Database `json:"databases,omitempty"`
}

func (o listDbSubscription) String() string {
	return internal.ToString(o)
}
