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
