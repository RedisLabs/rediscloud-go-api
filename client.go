package rediscloud_api

import (
	"log"
	"net/http"
	"os"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/service/account"
	"github.com/RedisLabs/rediscloud-go-api/service/cloud_accounts"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
	"github.com/RedisLabs/rediscloud-go-api/service/task"
)

type Client struct {
	Account      *account.API
	CloudAccount *cloud_accounts.API
	Database     *databases.API
	Subscription *subscriptions.API
	Task         *task.API
}

func NewClient(configs ...Option) (*Client, error) {
	config := &Options{
		baseUrl:   "https://api.redislabs.com/v1",
		userAgent: userAgent,
		apiKey:    os.Getenv(AccessKeyEnvVar),
		secretKey: os.Getenv(SecretKeyEnvVar),
		logger:    log.New(os.Stderr, "", log.LstdFlags),
		transport: http.DefaultTransport,
	}

	for _, option := range configs {
		option(config)
	}

	httpClient := &http.Client{
		Transport: config.roundTripper(),
	}

	client, err := internal.NewHttpClient(httpClient, config.baseUrl)
	if err != nil {
		return nil, err
	}

	t := task.NewAPI(client, config.logger)

	a := account.NewAPI(client)
	c := cloud_accounts.NewAPI(client, t, config.logger)
	d := databases.NewAPI(client, t, config.logger)
	s := subscriptions.NewAPI(client, t, config.logger)

	return &Client{
		Account:      a,
		CloudAccount: c,
		Database:     d,
		Subscription: s,
		Task:         t,
	}, nil
}

type Options struct {
	baseUrl   string
	apiKey    string
	secretKey string
	userAgent string
	logger    Log
	transport http.RoundTripper
}

func (o Options) roundTripper() http.RoundTripper {
	return &credentialTripper{
		apiKey:    o.apiKey,
		secretKey: o.secretKey,
		wrapped:   o.transport,
	}
}

type Option func(*Options)

func Auth(apiKey string, secretKey string) Option {
	return func(options *Options) {
		options.apiKey = apiKey
		options.secretKey = secretKey
	}
}

func BaseUrl(url string) Option {
	return func(options *Options) {
		options.baseUrl = url
	}
}

func Transporter(transporter http.RoundTripper) Option {
	return func(options *Options) {
		options.transport = transporter
	}
}

func AdditionalUserAgent(additional string) Option {
	return func(options *Options) {
		options.userAgent += " " + additional
	}
}

func Logger(log Log) Option {
	return func(options *Options) {
		options.logger = log
	}
}

type Log interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type credentialTripper struct {
	apiKey    string
	secretKey string
	wrapped   http.RoundTripper
}

func (c *credentialTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("X-Api-Secret-Key", c.secretKey)

	return c.wrapped.RoundTrip(request)
}

var _ http.RoundTripper = &credentialTripper{}
