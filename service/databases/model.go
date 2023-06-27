package databases

import (
	"fmt"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type CreateDatabase struct {
	DryRun                              *bool                        `json:"dryRun,omitempty"`
	Name                                *string                      `json:"name,omitempty"`
	Protocol                            *string                      `json:"protocol,omitempty"`
	MemoryLimitInGB                     *float64                     `json:"memoryLimitInGb,omitempty"`
	SupportOSSClusterAPI                *bool                        `json:"supportOSSClusterApi,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool                        `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	DataPersistence                     *string                      `json:"dataPersistence,omitempty"`
	DataEvictionPolicy                  *string                      `json:"dataEvictionPolicy,omitempty"`
	Replication                         *bool                        `json:"replication,omitempty"`
	ThroughputMeasurement               *CreateThroughputMeasurement `json:"throughputMeasurement,omitempty"`
	AverageItemSizeInBytes              *int                         `json:"averageItemSizeInBytes,omitempty"`
	ReplicaOf                           []*string                    `json:"replicaOf,omitempty"`
	// Deprecated: Use RemoteBackup instead
	PeriodicBackupPath   *string               `json:"periodicBackupPath,omitempty"`
	SourceIP             []*string             `json:"sourceIp,omitempty"`
	ClientSSLCertificate *string               `json:"clientSslCertificate,omitempty"`
	Password             *string               `json:"password,omitempty"`
	Alerts               []*CreateAlert        `json:"alerts,omitempty"`
	Modules              []*CreateModule       `json:"modules,omitempty"`
	EnableTls            *bool                 `json:"enableTls,omitempty"`
	PortNumber           *int                  `json:"port,omitempty"`
	RemoteBackup         *DatabaseBackupConfig `json:"remoteBackup,omitempty"`
}

func (o CreateDatabase) String() string {
	return internal.ToString(o)
}

type CreateThroughputMeasurement struct {
	By    *string `json:"by,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o CreateThroughputMeasurement) String() string {
	return internal.ToString(o)
}

type CreateAlert struct {
	Name  *string `json:"name,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o CreateAlert) String() string {
	return internal.ToString(o)
}

type CreateModule struct {
	Name *string `json:"name,omitempty"`
}

func (o CreateModule) String() string {
	return internal.ToString(o)
}

// TODO: do we need a separate ActiveActiveDatabase type as well?
type Database struct {
	ID                     *int        `json:"databaseId,omitempty"`
	Name                   *string     `json:"name,omitempty"`
	Protocol               *string     `json:"protocol,omitempty"`
	Provider               *string     `json:"provider,omitempty"`
	Region                 *string     `json:"region,omitempty"`
	Status                 *string     `json:"status,omitempty"`
	MemoryLimitInGB        *float64    `json:"memoryLimitInGb,omitempty"`
	MemoryUsedInMB         *float64    `json:"memoryUsedInMb,omitempty"`
	SupportOSSClusterAPI   *bool       `json:"supportOSSClusterApi,omitempty"`
	DataPersistence        *string     `json:"dataPersistence,omitempty"`
	Replication            *bool       `json:"replication,omitempty"`
	DataEvictionPolicy     *string     `json:"dataEvictionPolicy,omitempty"`
	ThroughputMeasurement  *Throughput `json:"throughputMeasurement,omitempty"`
	ReplicaOf              *ReplicaOf  `json:"replicaOf,omitempty"`
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
	Backup                 *Backup     `json:"backup,omitempty"`
}

func (o Database) String() string {
	return internal.ToString(o)
}

type ReplicaOf struct {
	Endpoints []*string `json:"endpoints,omitempty"`
}

type Clustering struct {
	NumberOfShards *int         `json:"numberOfShards,omitempty"`
	RegexRules     []*RegexRule `json:"regexRules,omitempty"`
	// TODO HashingPolicy interface{} `json:"hashingPolicy,omitempty"`
}

func (o Clustering) String() string {
	return internal.ToString(o)
}

type RegexRule struct {
	Ordinal int    `json:"ordinal"`
	Pattern string `json:"pattern"`
}

func (o RegexRule) String() string {
	return internal.ToString(o)
}

type Security struct {
	SSLClientAuthentication *bool     `json:"sslClientAuthentication,omitempty"`
	SourceIPs               []*string `json:"sourceIps,omitempty"`
	Password                *string   `json:"password,omitempty"`
	EnableTls               *bool     `json:"enableTls,omitempty"`
}

func (o Security) String() string {
	return internal.ToString(o)
}

type Module struct {
	Name *string `json:"name,omitempty"`
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

type UpdateDatabase struct {
	DryRun                              *bool                        `json:"dryRun,omitempty"`
	Name                                *string                      `json:"name,omitempty"`
	MemoryLimitInGB                     *float64                     `json:"memoryLimitInGb,omitempty"`
	SupportOSSClusterAPI                *bool                        `json:"supportOSSClusterApi,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool                        `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	DataEvictionPolicy                  *string                      `json:"dataEvictionPolicy,omitempty"`
	Replication                         *bool                        `json:"replication,omitempty"`
	ThroughputMeasurement               *UpdateThroughputMeasurement `json:"throughputMeasurement,omitempty"`
	RegexRules                          []*string                    `json:"regexRules,omitempty"`
	DataPersistence                     *string                      `json:"dataPersistence,omitempty"`
	ReplicaOf                           []*string                    `json:"replicaOf"`
	PeriodicBackupPath                  *string                      `json:"periodicBackupPath,omitempty"`
	SourceIP                            []*string                    `json:"sourceIp,omitempty"`
	ClientSSLCertificate                *string                      `json:"clientSslCertificate,omitempty"`
	Password                            *string                      `json:"password,omitempty"`
	Alerts                              []*UpdateAlert               `json:"alerts,omitempty"`
	EnableTls                           *bool                        `json:"enableTls,omitempty"`
	RemoteBackup                        *DatabaseBackupConfig        `json:"remoteBackup,omitempty"`
}

func (o UpdateDatabase) String() string {
	return internal.ToString(o)
}

type UpdateThroughputMeasurement struct {
	By    *string `json:"by,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o UpdateThroughputMeasurement) String() string {
	return internal.ToString(o)
}

type UpdateAlert struct {
	Name  *string `json:"name,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o UpdateAlert) String() string {
	return internal.ToString(o)
}

type Import struct {
	SourceType    *string   `json:"sourceType,omitempty"`
	ImportFromURI []*string `json:"importFromUri,omitempty"`
}

func (o Import) String() string {
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

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("database %d in subscription %d not found", f.dbId, f.subId)
}

const (
	// StatusActive is the active value of the `Status` field in `Database`
	StatusActive = "active"
	// StatusDraft is the draft value of the `Status` field in `Database`
	StatusDraft = "draft"
	// StatusPending is the pending value of the `Status` field in `Database`
	StatusPending = "pending"
	// StatusRCPChangePending is the RCP change pending value of the `Status` field in `Database`
	StatusRCPChangePending = "rcp-change-pending"
	// StatusRCPDraft is the RCP draft value of the `Status` field in `Database`
	StatusRCPDraft = "rcp-draft"
	// StatusRCPActiveChangeDraft is the RCP active change draft value of the `Status` field in `Database`
	StatusRCPActiveChangeDraft = "rcp-active-change-draft"
	// StatusActiveChangeDraft is the Active change draft value of the `Status` field in `Database`
	StatusActiveChangeDraft = "active-change-draft"
	// StatusActiveChangePending is the Active change pending value of the `Status` field in `Database`
	StatusActiveChangePending = "active-change-pending"

	// StatusProxyPolicyChangePending and StatusProxyPolicyChangeDraft
	//The below two Proxy Policy states are caused by a change to the 'support_oss_cluster_api' attribute
	// StatusProxyPolicyChangePending is the Proxy Policy change pending value of the `Status` field in `Database`.
	StatusProxyPolicyChangePending = "proxy-policy-change-pending"
	// StatusProxyPolicyChangeDraft is the Proxy Policy change draft value of the `Status` field in `Database`
	StatusProxyPolicyChangeDraft = "proxy-policy-change-draft"

	// StatusError is the error value of the `Status` field in `Database`
	StatusError = "error"
	// BackupIntervalEvery24Hours is the schedule to back up once a day
	BackupIntervalEvery24Hours = "every-24-hours"
	// BackupIntervalEvery12Hours is the schedule to back up twice a day
	BackupIntervalEvery12Hours = "every-12-hours"
	// BackupIntervalEvery6Hours is the schedule to back up four times a day
	BackupIntervalEvery6Hours = "every-6-hours"
	// BackupIntervalEvery4Hours is the schedule to back up six times a day
	BackupIntervalEvery4Hours = "every-4-hours"
	// BackupIntervalEvery2Hours is the schedule to back up twelve times a day
	BackupIntervalEvery2Hours = "every-2-hours"
	// BackupIntervalEvery1Hours is the schedule to back up every hour
	BackupIntervalEvery1Hours = "every-1-hours"
	// MemoryStorageRam stores data only in RAM
	MemoryStorageRam = "ram"
	// MemoryStorageRamAndFlash stores data both in RAM and on SSD
	MemoryStorageRamAndFlash = "ram-and-flash"
)

func MemoryStorageValues() []string {
	return []string{
		MemoryStorageRam,
		MemoryStorageRamAndFlash,
	}
}

func ProtocolValues() []string {
	return []string{
		"redis",
		"memcached",
	}
}

func DataPersistenceValues() []string {
	return []string{
		"none",
		"aof-every-1-second",
		"aof-every-write",
		"snapshot-every-1-hour",
		"snapshot-every-6-hours",
		"snapshot-every-12-hours",
	}
}

func DataEvictionPolicyValues() []string {
	return []string{
		"allkeys-lru",
		"allkeys-lfu",
		"allkeys-random",
		"volatile-lru",
		"volatile-lfu",
		"volatile-random",
		"volatile-ttl",
		"noeviction",
	}
}

func SourceTypeValues() []string {
	return []string{
		"http",
		"redis",
		"ftp",
		"aws-s3",
		"azure-blob-storage",
		"google-blob-storage",
	}
}

func AlertNameValues() []string {
	return []string{
		"dataset-size",
		"throughput-higher-than",
		"throughput-lower-than",
		"latency",
		"syncsource-error",
		"syncsource-lag",
	}
}

func BackupStorageTypes() []string {
	return []string{
		"ftp",
		"aws-s3",
		"azure-blob-storage",
		"google-blob-storage",
	}
}

func BackupIntervals() []string {
	return []string{
		BackupIntervalEvery24Hours,
		BackupIntervalEvery12Hours,
		BackupIntervalEvery6Hours,
		BackupIntervalEvery4Hours,
		BackupIntervalEvery2Hours,
		BackupIntervalEvery1Hours,
	}
}
