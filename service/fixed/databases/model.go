package databases

import (
	"fmt"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
)

type CreateFixedDatabase struct {
	Name                                *string                `json:"name,omitempty"`
	Protocol                            *string                `json:"protocol,omitempty"`
	MemoryLimitInGB                     *float64               `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB                     *float64               `json:"datasetSizeInGb,omitempty"`
	SupportOSSClusterAPI                *bool                  `json:"supportOSSClusterApi,omitempty"`
	RespVersion                         *string                `json:"respVersion,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool                  `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	EnableDatabaseClustering            *bool                  `json:"enableDatabaseClustering,omitempty"`
	DataPersistence                     *string                `json:"dataPersistence,omitempty"`
	DataEvictionPolicy                  *string                `json:"dataEvictionPolicy,omitempty"`
	Replication                         *bool                  `json:"replication,omitempty"`
	PeriodicBackupPath                  *string                `json:"periodicBackupPath,omitempty"`
	SourceIPs                           []*string              `json:"sourceIps,omitempty"`
	RegexRules                          []*string              `json:"regexRules,omitempty"`
	Replica                             *ReplicaOf             `json:"replica,omitempty"`
	ClientTlsCertificates               []*DatabaseCertificate `json:"clientTlsCertificates,omitempty"`
	EnableTls                           *bool                  `json:"enableTls,omitempty"`
	Password                            *string                `json:"password,omitempty"`
	Alerts                              *[]*databases.Alert    `json:"alerts,omitempty"`
	Modules                             *[]*databases.Module   `json:"modules,omitempty"`
}

type UpdateFixedDatabase struct {
	Name                                *string                `json:"name,omitempty"`
	MemoryLimitInGB                     *float64               `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB                     *float64               `json:"datasetSizeInGb,omitempty"`
	SupportOSSClusterAPI                *bool                  `json:"supportOSSClusterApi,omitempty"`
	RespVersion                         *string                `json:"respVersion,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool                  `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	EnableDatabaseClustering            *bool                  `json:"enableDatabaseClustering,omitempty"`
	DataPersistence                     *string                `json:"dataPersistence,omitempty"`
	DataEvictionPolicy                  *string                `json:"dataEvictionPolicy,omitempty"`
	Replication                         *bool                  `json:"replication,omitempty"`
	PeriodicBackupPath                  *string                `json:"periodicBackupPath,omitempty"`
	SourceIPs                           []*string              `json:"sourceIps,omitempty"`
	Replica                             *ReplicaOf             `json:"replica,omitempty"`
	RegexRules                          []*string              `json:"regexRules,omitempty"`
	ClientTlsCertificates               []*DatabaseCertificate `json:"clientTlsCertificates,omitempty"`
	EnableTls                           *bool                  `json:"enableTls,omitempty"`
	Password                            *string                `json:"password,omitempty"`
	Alerts                              *[]*databases.Alert    `json:"alerts,omitempty"`
	// As with flexible databases, this is only available on the update endpoint
	EnableDefaultUser *bool `json:"enableDefaultUser,omitempty"`
}

type FixedDatabase struct {
	DatabaseId                          *int       `json:"databaseId,omitempty"`
	Name                                *string    `json:"name,omitempty"`
	Protocol                            *string    `json:"protocol,omitempty"`
	Provider                            *string    `json:"provider,omitempty"`
	Region                              *string    `json:"region,omitempty"`
	RedisVersionCompliance              *string    `json:"redisVersionCompliance,omitempty"`
	RespVersion                         *string    `json:"respVersion,omitempty"`
	Status                              *string    `json:"status,omitempty"`
	PlanMemoryLimit                     *float64   `json:"planMemoryLimit,omitempty"`
	MemoryLimitMeasurementUnit          *string    `json:"memoryLimitMeasurementUnit,omitempty"`
	MemoryLimitInGb                     *float64   `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB                     *float64   `json:"datasetSizeInGb,omitempty"`
	MemoryUsedInMb                      *float64   `json:"memoryUsedInMb,omitempty"`
	NetworkMonthlyUsageInByte           *float64   `json:"networkMonthlyUsageInByte,omitempty"`
	MemoryStorage                       *string    `json:"memoryStorage,omitempty"`
	SupportOSSClusterAPI                *bool      `json:"supportOSSClusterApi,omitempty"`
	UseExternalEndpointForOSSClusterAPI *bool      `json:"useExternalEndpointForOSSClusterApi,omitempty"`
	DataPersistence                     *string    `json:"dataPersistence,omitempty"`
	Replication                         *bool      `json:"replication,omitempty"`
	DataEvictionPolicy                  *string    `json:"dataEvictionPolicy,omitempty"`
	ActivatedOn                         *time.Time `json:"activatedOn,omitempty"`
	LastModified                        *time.Time `json:"lastModified,omitempty"`
	PublicEndpoint                      *string    `json:"publicEndpoint,omitempty"`
	PrivateEndpoint                     *string    `json:"privateEndpoint,omitempty"`
	// The following are undocumented but are returned
	Replica    *ReplicaOf           `json:"replica,omitempty"`
	Clustering *Clustering          `json:"clustering,omitempty"`
	Security   *Security            `json:"security,omitempty"`
	Modules    *[]*databases.Module `json:"modules,omitempty"`
	Alerts     *[]*databases.Alert  `json:"alerts,omitempty"`
	Backup     *Backup              `json:"backup,omitempty"`
}

type ReplicaOf struct {
	Description *string       `json:"description,omitempty"`
	SyncSources []*SyncSource `json:"syncSources,omitempty"`
}

func (o ReplicaOf) String() string {
	return internal.ToString(o)
}

type SyncSource struct {
	Description *string `json:"description,omitempty"`
	Endpoint    *string `json:"endpoint,omitempty"`
	Encryption  *bool   `json:"encryption,omitempty"`
	ServerCert  *string `json:"serverCert,omitempty"`
}

func (o SyncSource) String() string {
	return internal.ToString(o)
}

type DatabaseCertificate struct {
	Description                *string `json:"description,omitempty"`
	PublicCertificatePEMString *string `json:"publicCertificatePEMString,omitempty"`
}

func (o DatabaseCertificate) String() string {
	return internal.ToString(o)
}

type Clustering struct {
	Enabled       *bool                  `json:"enabled,omitempty"`
	RegexRules    []*databases.RegexRule `json:"regexRules,omitempty"`
	HashingPolicy *string                `json:"hashingPolicy,omitempty"`
}

func (o Clustering) String() string {
	return internal.ToString(o)
}

type Security struct {
	EnableDefaultUser       *bool     `json:"defaultUserEnabled,omitempty"`
	Password                *string   `json:"password,omitempty"`
	SSLClientAuthentication *bool     `json:"sslClientAuthentication,omitempty"`
	TLSClientAuthentication *bool     `json:"tlsClientAuthentication,omitempty"`
	EnableTls               *bool     `json:"enableTls,omitempty"`
	SourceIPs               []*string `json:"sourceIps,omitempty"`
}

func (o Security) String() string {
	return internal.ToString(o)
}

type Backup struct {
	Enabled     *bool   `json:"remoteBackupEnabled,omitempty"`
	Status      *string `json:"status,omitempty"`
	Interval    *string `json:"interval,omitempty"`
	Destination *string `json:"destination,omitempty"`
}

func (o Backup) String() string {
	return internal.ToString(o)
}

type Import struct {
	SourceType    *string   `json:"sourceType,omitempty"`
	ImportFromURI []*string `json:"importFromUri,omitempty"`
}

func (o Import) String() string {
	return internal.ToString(o)
}

type listFixedDatabaseResponse struct {
	FixedSubscription *listDbSubscription `json:"subscription,omitempty"`
}

func (o listFixedDatabaseResponse) String() string {
	return internal.ToString(o)
}

type listDbSubscription struct {
	ID        *int             `json:"subscriptionId,omitempty"`
	Databases []*FixedDatabase `json:"databases,omitempty"`
}

func (o listDbSubscription) String() string {
	return internal.ToString(o)
}

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("fixed database %d in subscription %d not found", f.dbId, f.subId)
}

func ProtocolValues() []string {
	return []string{
		"redis",
		"memcached",
		"stack",
	}
}
