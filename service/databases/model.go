package databases

import "time"

type taskResponse struct {
	TaskId string `json:"taskId"`
}

type Database struct {
	ID                   int     `json:"databaseId"`
	Name                 string  `json:"name"`
	Protocol             string  `json:"protocol"`
	Provider             string  `json:"provider"`
	Region               string  `json:"region"`
	Status               string  `json:"status"`
	MemoryLimitInGb      float64 `json:"memoryLimitInGb,omitempty"`
	MemoryUsedInMb       float64 `json:"memoryUsedInMb"`
	SupportOSSClusterApi bool    `json:"supportOSSClusterApi"`
	DataPersistence      string  `json:"dataPersistence"`
	Replication          bool    `json:"replication"`
	DataEvictionPolicy   string  `json:"dataEvictionPolicy"`
	// TODO throughputMeasurement
	// TODO replicaOf
	// TODO clustering
	// TODO security
	// TODO modules
	// TODO alerts
	ActivatedOn            time.Time `json:"activatedON"`
	LastModified           time.Time `json:"lastModified"`
	MemoryStorage          string    `json:"memoryStorage,omitempty"`
	PrivateEndpoint        string    `json:"privateEndpoint"`
	PublicEndpoint         string    `json:"publicEndpoint"`
	RedisVersionCompliance string    `json:"redisVersionCompliance"`
}

type listDatabaseResponse struct {
	Subscription []listDbSubscription `json:"subscription"`
}

type listDbSubscription struct {
	ID        int         `json:"subscriptionId"`
	Databases []*Database `json:"databases"`
}
