package cloud_accounts

import (
	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type CreateCloudAccount struct {
	AccessKeyID     *string `json:"accessKeyId,omitempty"`
	AccessSecretKey *string `json:"accessSecretKey,omitempty"`
	ConsoleUsername *string `json:"consoleUsername,omitempty"`
	ConsolePassword *string `json:"consolePassword,omitempty"`
	Name            *string `json:"name,omitempty"`
	Provider        *string `json:"provider,omitempty"`
	SignInLoginURL  *string `json:"signInLoginUrl,omitempty"`
}

func (o CreateCloudAccount) String() string {
	return internal.ToString(o)
}

type UpdateCloudAccount struct {
	AccessKeyID     *string `json:"accessKeyId,omitempty"`
	AccessSecretKey *string `json:"accessSecretKey,omitempty"`
	ConsoleUsername *string `json:"consoleUsername,omitempty"`
	ConsolePassword *string `json:"consolePassword,omitempty"`
	Name            *string `json:"name,omitempty"`
	SignInLoginURL  *string `json:"signInLoginUrl,omitempty"`
}

func (o UpdateCloudAccount) String() string {
	return internal.ToString(o)
}

type listCloudAccounts struct {
	CloudAccounts []*CloudAccount `json:"cloudAccounts"`
}

type CloudAccount struct {
	ID          *int    `json:"id"`
	Name        *string `json:"name,omitempty"`
	Provider    *string `json:"provider,omitempty"`
	Status      *string `json:"status,omitempty"`
	AccessKeyID *string `json:"accessKeyId,omitempty"`
}

func (o CloudAccount) String() string {
	return internal.ToString(o)
}
