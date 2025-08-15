package rediscloud_api

import (
	"testing"

	fixedDatabases "github.com/RedisLabs/rediscloud-go-api/service/fixed/databases"

	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/redis_rules"
	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/roles"
	"github.com/RedisLabs/rediscloud-go-api/service/cloud_accounts"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	fixedSubscriptions "github.com/RedisLabs/rediscloud-go-api/service/fixed/subscriptions"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
	"github.com/stretchr/testify/assert"
)

func TestSubcriptionFixtures(t *testing.T) {
	assert.Equal(t, "active", subscriptions.SubscriptionStatusActive)
	assert.Equal(t, "pending", subscriptions.SubscriptionStatusPending)
	assert.Equal(t, "error", subscriptions.SubscriptionStatusError)
	assert.Equal(t, "deleting", subscriptions.SubscriptionStatusDeleting)

	assert.Equal(t, "initiating-request", subscriptions.VPCPeeringStatusInitiatingRequest)
	assert.Equal(t, "active", subscriptions.VPCPeeringStatusActive)
	assert.Equal(t, "inactive", subscriptions.VPCPeeringStatusInactive)
	assert.Equal(t, "pending-acceptance", subscriptions.VPCPeeringStatusPendingAcceptance)
	assert.Equal(t, "failed", subscriptions.VPCPeeringStatusFailed)

	assert.Equal(t, "single-region", subscriptions.SubscriptionDeploymentTypeSingleRegion)
	assert.Equal(t, "active-active", subscriptions.SubscriptionDeploymentTypeActiveActive)
}

func TestDatabaseFixtures(t *testing.T) {
	assert.Equal(t, "active", databases.StatusActive)
	assert.Equal(t, "draft", databases.StatusDraft)
	assert.Equal(t, "pending", databases.StatusPending)
	assert.Equal(t, "rcp-change-pending", databases.StatusRCPChangePending)
	assert.Equal(t, "rcp-draft", databases.StatusRCPDraft)
	assert.Equal(t, "rcp-active-change-draft", databases.StatusRCPActiveChangeDraft)
	assert.Equal(t, "active-change-draft", databases.StatusActiveChangeDraft)
	assert.Equal(t, "active-change-pending", databases.StatusActiveChangePending)
	assert.Equal(t, "dynamic-endpoints-creation-pending", databases.StatusDynamicEndpointsCreationPending)
	assert.Equal(t, "active-upgrade-pending", databases.StatusActiveUpgradePending)

	assert.Equal(t, "proxy-policy-change-pending", databases.StatusProxyPolicyChangePending)
	assert.Equal(t, "proxy-policy-change-draft", databases.StatusProxyPolicyChangeDraft)
	assert.Equal(t, "error", databases.StatusError)

	assert.Equal(t, []string{databases.MemoryStorageRam, databases.MemoryStorageRamAndFlash}, databases.MemoryStorageValues())
	assert.Equal(t, []string{"redis", "memcached"}, databases.ProtocolValues())
	assert.Equal(t, []string{
		"none",
		"aof-every-1-second",
		"aof-every-write",
		"snapshot-every-1-hour",
		"snapshot-every-6-hours",
		"snapshot-every-12-hours",
	}, databases.DataPersistenceValues())
	assert.Equal(t, []string{
		"allkeys-lru",
		"allkeys-lfu",
		"allkeys-random",
		"volatile-lru",
		"volatile-lfu",
		"volatile-random",
		"volatile-ttl",
		"noeviction",
	}, databases.DataEvictionPolicyValues())
	assert.Equal(t, []string{
		"http",
		"redis",
		"ftp",
		"aws-s3",
		"azure-blob-storage",
		"google-blob-storage",
	}, databases.SourceTypeValues())
	assert.Equal(t, []string{
		"dataset-size",
		"throughput-higher-than",
		"throughput-lower-than",
		"latency",
		"syncsource-error",
		"syncsource-lag",
	}, databases.AlertNameValues())
	assert.Equal(t, []string{
		"ftp",
		"aws-s3",
		"azure-blob-storage",
		"google-blob-storage",
	}, databases.BackupStorageTypes())
	assert.Equal(t, []string{
		databases.BackupIntervalEvery24Hours,
		databases.BackupIntervalEvery12Hours,
		databases.BackupIntervalEvery6Hours,
		databases.BackupIntervalEvery4Hours,
		databases.BackupIntervalEvery2Hours,
		databases.BackupIntervalEvery1Hours,
	}, databases.BackupIntervals())
}

func TestCloudAccountFixtures(t *testing.T) {
	assert.Equal(t, "active", cloud_accounts.StatusActive)
	assert.Equal(t, "draft", cloud_accounts.StatusDraft)
	assert.Equal(t, "change-draft", cloud_accounts.StatusChangeDraft)
	assert.Equal(t, "error", cloud_accounts.StatusError)
	assert.Equal(t, []string{"AWS", "GCP"}, cloud_accounts.ProviderValues())
}

func TestRedisRuleFixtures(t *testing.T) {
	assert.Equal(t, "active", redis_rules.StatusActive)
	assert.Equal(t, "pending", redis_rules.StatusPending)
	assert.Equal(t, "error", redis_rules.StatusError)
	assert.Equal(t, "deleting", redis_rules.StatusDeleting)
}

func TestRoleFixtures(t *testing.T) {
	assert.Equal(t, "active", roles.StatusActive)
	assert.Equal(t, "pending", roles.StatusPending)
	assert.Equal(t, "error", roles.StatusError)
	assert.Equal(t, "deleting", roles.StatusDeleting)
}

func TestFixedSubcriptionFixtures(t *testing.T) {
	assert.Equal(t, "active", fixedSubscriptions.FixedSubscriptionStatusActive)
	assert.Equal(t, "pending", fixedSubscriptions.FixedSubscriptionStatusPending)
	assert.Equal(t, "error", fixedSubscriptions.FixedSubscriptionStatusError)
	assert.Equal(t, "deleting", fixedSubscriptions.FixedSubscriptionStatusDeleting)

	assert.Equal(t, []string{"redis", "memcached", "stack"}, fixedDatabases.ProtocolValues())
}
