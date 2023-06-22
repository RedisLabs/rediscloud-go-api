package users

import "github.com/RedisLabs/rediscloud-go-api/internal"

// Get

type ListUsersResponse struct {
	AccountId *int               `json:"accountId,omitempty"`
	Users     []*GetUserResponse `json:"users,omitempty"`
}

func (o ListUsersResponse) String() string {
	return internal.ToString(o)
}

type GetUserResponse struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Role *string `json:"role,omitempty"`
}

func (o GetUserResponse) String() string {
	return internal.ToString(o)
}

// Create

type CreateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Role     *string `json:"role,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (o CreateUserRequest) String() string {
	return internal.ToString(o)
}

// Update

type UpdateUserRequest struct {
	Role     *string `json:"role,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (o UpdateUserRequest) String() string {
	return internal.ToString(o)
}
