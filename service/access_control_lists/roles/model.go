package roles

import "github.com/RedisLabs/rediscloud-go-api/internal"

// Read

type ListRolesResponse struct {
	AccountId *int               `json:"accountId,omitempty"`
	Roles     []*GetRoleResponse `json:"roles,omitempty"`
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
	RuleId    *int                               `json:"ruleId,omitempty"`
	RuleName  *string                            `json:"ruleName,omitempty"`
	Databases []*GetDatabaseInRuleInRoleResponse `json:"databases,omitempty"`
}

func (o GetRuleInRoleResponse) String() string {
	return internal.ToString(o)
}

type GetDatabaseInRuleInRoleResponse struct {
	SubscriptionId *int      `json:"subscriptionId,omitempty"`
	DatabaseId     *int      `json:"databaseId,omitempty"`
	DatabaseName   *string   `json:"databaseName,omitempty"`
	Regions        []*string `json:"regions,omitempty"`
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
	Regions        []*string `json:"regions,omitempty"`
}

func (o CreateDatabaseInRuleInRoleRequest) String() string {
	return internal.ToString(o)
}

const (
	// StatusActive is the active value of the `Status` field in `Role`
	StatusActive = "active"
	// StatusPending is the pending value of the `Status` field in `Role`
	StatusPending = "pending"
	// StatusError is the error value of the `Status` field in `Role`
	StatusError = "error"
	// StatusDeleting is the deleting value of the `Status` field in `Role`
	StatusDeleting = "deleting"
)
