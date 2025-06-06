package databases

import (
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type ActiveActiveDatabase struct {
	ID                                  *int            `json:"databaseId,omitempty"`
	Name                                *string         `json:"name,omitempty"`
	Protocol                            *string         `json:"protocol,omitempty"`
	Status                              *string         `json:"status,omitempty"`
	MemoryStorage                       *string         `json:"memoryStorage,omitempty"`
	ActiveActiveRedis                   *bool           `json:"activeActiveRedis,omitempty"`
	ActivatedOn                         *time.Time      `json:"activatedOn,omitempty"`
	LastModified                        *time.Time      `json:"lastModified,omitempty"`
	SupportOSSClusterAPI                *bool           `json:"supportOSSClusterApi,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool           `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	Replication                         *bool           `json:"replication,omitempty"`
	DataEvictionPolicy                  *string         `json:"dataEvictionPolicy,omitempty"`
	Modules                             []*Module       `json:"modules,omitempty"`
	CrdbDatabases                       []*CrdbDatabase `json:"crdbDatabases,omitempty"`
}

func (o ActiveActiveDatabase) String() string {
	return internal.ToString(o)
}

type CrdbDatabase struct {
	Provider                 *string   `json:"provider,omitempty"`
	Region                   *string   `json:"region,omitempty"`
	RedisVersionCompliance   *string   `json:"redisVersionCompliance,omitempty"`
	PublicEndpoint           *string   `json:"publicEndpoint,omitempty"`
	PrivateEndpoint          *string   `json:"privateEndpoint,omitempty"`
	MemoryLimitInGB          *float64  `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB          *float64  `json:"datasetSizeInGb,omitempty"`
	MemoryUsedInMB           *float64  `json:"memoryUsedInMb,omitempty"`
	ReadOperationsPerSecond  *int      `json:"readOperationsPerSecond,omitempty"`
	WriteOperationsPerSecond *int      `json:"writeOperationsPerSecond,omitempty"`
	DataPersistence          *string   `json:"dataPersistence,omitempty"`
	Alerts                   []*Alert  `json:"alerts,omitempty"`
	Security                 *Security `json:"security,omitempty"`
	Backup                   *Backup   `json:"backup,omitempty"`
	QueryPerformanceFactor   *string   `json:"queryPerformanceFactor,omitempty"`
}

func (o CrdbDatabase) String() string {
	return internal.ToString(o)
}

type Backup struct {
	Enabled     *bool   `json:"enableRemoteBackup,omitempty"`
	TimeUTC     *string `json:"timeUTC,omitempty"`
	Interval    *string `json:"interval,omitempty"`
	Destination *string `json:"destination,omitempty"`
}

func (o Backup) String() string {
	return internal.ToString(o)
}

type CreateActiveActiveDatabase struct {
	DryRun                              *bool              `json:"dryRun,omitempty"`
	Name                                *string            `json:"name,omitempty"`
	Protocol                            *string            `json:"protocol,omitempty"`
	MemoryLimitInGB                     *float64           `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB                     *float64           `json:"datasetSizeInGb,omitempty"`
	SupportOSSClusterAPI                *bool              `json:"supportOSSClusterApi,omitempty"`
	RespVersion                         *string            `json:"respVersion,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool              `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	DataEvictionPolicy                  *string            `json:"dataEvictionPolicy,omitempty"`
	GlobalDataPersistence               *string            `json:"dataPersistence,omitempty"`
	GlobalSourceIP                      []*string          `json:"sourceIp,omitempty"`
	GlobalPassword                      *string            `json:"password,omitempty"`
	GlobalAlerts                        []*Alert           `json:"alerts,omitempty"`
	GlobalModules                       []*Module          `json:"modules,omitempty"`
	LocalThroughputMeasurement          []*LocalThroughput `json:"localThroughputMeasurement,omitempty"`
	PortNumber                          *int               `json:"port,omitempty"`
	QueryPerformanceFactor              *string            `json:"queryPerformanceFactor,omitempty"`
}

func (o CreateActiveActiveDatabase) String() string {
	return internal.ToString(o)
}

type LocalThroughput struct {
	Region                   *string `json:"region,omitempty"`
	WriteOperationsPerSecond *int    `json:"writeOperationsPerSecond,omitempty"`
	ReadOperationsPerSecond  *int    `json:"readOperationsPerSecond,omitempty"`
}

func (o LocalThroughput) String() string {
	return internal.ToString(o)
}

type UpdateActiveActiveDatabase struct {
	DryRun                              *bool    `json:"dryRun,omitempty"`
	MemoryLimitInGB                     *float64 `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB                     *float64 `json:"datasetSizeInGb,omitempty"`
	SupportOSSClusterAPI                *bool    `json:"supportOSSClusterApi,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool    `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	ClientSSLCertificate                *string  `json:"clientSslCertificate,omitempty"`
	// Using a pointer to allow empty slices to be serialised/sent
	ClientTLSCertificates *[]*string `json:"clientTlsCertificates,omitempty"`
	EnableTls             *bool      `json:"enableTls,omitempty"`
	GlobalDataPersistence *string    `json:"globalDataPersistence,omitempty"`
	GlobalPassword        *string    `json:"globalPassword,omitempty"`
	GlobalSourceIP        []*string  `json:"globalSourceIp,omitempty"`
	// Using a pointer to allow empty slices to be serialised/sent
	GlobalAlerts           *[]*Alert                `json:"globalAlerts,omitempty"`
	Regions                []*LocalRegionProperties `json:"regions,omitempty"`
	DataEvictionPolicy     *string                  `json:"dataEvictionPolicy,omitempty"`
	QueryPerformanceFactor *string                  `json:"queryPerformanceFactor,omitempty"`
}

func (o UpdateActiveActiveDatabase) String() string {
	return internal.ToString(o)
}

type LocalRegionProperties struct {
	Region                     *string               `json:"region,omitempty"`
	RemoteBackup               *DatabaseBackupConfig `json:"remoteBackup,omitempty"`
	LocalThroughputMeasurement *LocalThroughput      `json:"localThroughputMeasurement,omitempty"`
	DataPersistence            *string               `json:"dataPersistence,omitempty"`
	Password                   *string               `json:"password,omitempty"`
	SourceIP                   []*string             `json:"sourceIp,omitempty"`
	// Using a pointer to allow empty slices to be serialised/sent
	Alerts *[]*Alert `json:"alerts,omitempty"`
}

func (o LocalRegionProperties) String() string {
	return internal.ToString(o)
}

type DatabaseBackupConfig struct {
	Active      *bool   `json:"active,omitempty"`
	Interval    *string `json:"interval,omitempty"`
	TimeUTC     *string `json:"timeUTC,omitempty"`
	StorageType *string `json:"storageType,omitempty"`
	StoragePath *string `json:"storagePath,omitempty"`
}

func (o DatabaseBackupConfig) String() string {
	return internal.ToString(o)
}

type listActiveActiveDatabaseResponse struct {
	AccountId    *int                              `json:"accountId,omitempty"`
	Subscription []*listActiveActiveDbSubscription `json:"subscription,omitempty"`
}

func (o listActiveActiveDatabaseResponse) String() string {
	return internal.ToString(o)
}

type listActiveActiveDbSubscription struct {
	ID        *int                    `json:"subscriptionId,omitempty"`
	Databases []*ActiveActiveDatabase `json:"databases,omitempty"`
}

func (o listActiveActiveDbSubscription) String() string {
	return internal.ToString(o)
}
