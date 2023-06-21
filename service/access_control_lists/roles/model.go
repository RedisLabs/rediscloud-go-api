package roles

import "github.com/RedisLabs/rediscloud-go-api/internal"

// Read

type ListRolesResponse struct {
	AccountId *int               `json:"accountId,omitempty"`
	Roles     []*GetRoleResponse `json:"roles,omitempty"`
	// Links []*AnotherThing
}

type GetRoleResponse struct {
	ID         *int                     `json:"id,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	RedisRules []*GetRuleInRoleResponse `json:"redisRules,omitempty"`
	Users      []*GetUserInRoleResponse `json:"users,omitempty"`
	Status     *string                  `json:"status,omitempty"`
}

func (o GetRoleResponse) String() string {
	return internal.ToString(o)
}

type GetRuleInRoleResponse struct {
	RuleId    *int
	RuleName  *string
	Databases []*GetDatabaseInRuleInRoleResponse
}

func (o GetRuleInRoleResponse) String() string {
	return internal.ToString(o)
}

type GetDatabaseInRuleInRoleResponse struct {
	SubscriptionId *int      `json:"subscriptionId,omitempty"`
	DatabaseId     *int      `json:"databaseId,omitempty"`
	DatabaseName   *string   `json:"databaseName,omitempty"`
	Regions        []*string `json:"regions,omitempty"` // Docs are unclear
}

func (o GetDatabaseInRuleInRoleResponse) String() string {
	return internal.ToString(o)
}

type GetUserInRoleResponse struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// Create and Update

type CreateRoleRequest struct {
	Name       *string                    `json:"name,omitempty"`
	RedisRules []*CreateRuleInRoleRequest `json:"redisRules,omitempty"`
}

func (o CreateRoleRequest) String() string {
	return internal.ToString(o)
}

type CreateRuleInRoleRequest struct {
	RuleName  *string                              `json:"ruleName,omitempty"`
	Databases []*CreateDatabaseInRuleInRoleRequest `json:"databases,omitempty"`
}

func (o CreateRuleInRoleRequest) String() string {
	return internal.ToString(o)
}

type CreateDatabaseInRuleInRoleRequest struct {
	SubscriptionId *int      `json:"subscriptionId,omitempty"`
	DatabaseId     *int      `json:"databaseId,omitempty"`
	Regions        []*string `json:"regions,omitempty"` // Docs are unclear
}

func (o CreateDatabaseInRuleInRoleRequest) String() string {
	return internal.ToString(o)
}
