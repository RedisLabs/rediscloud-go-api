package redis_rules

import "github.com/RedisLabs/rediscloud-go-api/internal"

// Read

type ListRedisRulesResponse struct {
	AccountId  *int                    `json:"accountId,omitempty"`
	RedisRules []*GetRedisRuleResponse `json:"redisRules,omitempty"`
}

func (o ListRedisRulesResponse) String() string {
	return internal.ToString(o)
}

type GetRedisRuleResponse struct {
	ID        *int    `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	ACL       *string `json:"acl,omitempty"`
	IsDefault *bool   `json:"isDefault,omitempty"`
	Status    *string `json:"status,omitempty"`
}

func (o GetRedisRuleResponse) String() string {
	return internal.ToString(o)
}

// Create + Update

type CreateRedisRuleRequest struct {
	Name      *string `json:"name,omitempty"`
	RedisRule *string `json:"redisRule,omitempty"`
}

func (o CreateRedisRuleRequest) String() string {
	return internal.ToString(o)
}

type UpdateRedisRuleRequest struct {
	RedisRule *string `json:"redisRule,omitempty"`
}

func (o UpdateRedisRuleRequest) String() string {
	return internal.ToString(o)
}

const (
	// StatusActive is the active value of the `Status` field in `RedisRule`
	StatusActive = "active"
	// StatusPending is the pending value of the `Status` field in `RedisRule`
	StatusPending = "pending"
	// StatusError is the error value of the `Status` field in `RedisRule`
	StatusError = "error"
	// StatusDeleting is the deleting value of the `Status` field in `RedisRule`
	StatusDeleting = "deleting"
)
