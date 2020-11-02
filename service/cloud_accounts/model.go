package cloud_accounts

type CreateCloudAccount struct {
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	AccessSecretKey string `json:"accessSecretKey,omitempty"`
	ConsoleUsername string `json:"consoleUsername,omitempty"`
	ConsolePassword string `json:"consolePassword,omitempty"`
	Name            string `json:"name,omitempty"`
	Provider        string `json:"provider,omitempty"`
	SignInLoginUrl  string `json:"signInLoginUrl,omitempty"`
}

type UpdateCloudAccount struct {
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	AccessSecretKey string `json:"accessSecretKey,omitempty"`
	ConsoleUsername string `json:"consoleUsername,omitempty"`
	ConsolePassword string `json:"consolePassword,omitempty"`
	Name            string `json:"name,omitempty"`
	SignInLoginUrl  string `json:"signInLoginUrl,omitempty"`
}

type taskResponse struct {
	TaskId string `json:"taskId"`
}

type CloudAccount struct {
	Name        string `json:"name,omitempty"`
	Provider    string `json:"provider,omitempty"`
	Status      string `json:"status,omitempty"`
	AccessKeyId string `json:"accessKeyId,omitempty"`
}
