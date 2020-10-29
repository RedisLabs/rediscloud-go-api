package rediscloud_go_api

import (
	"net/http"
	"os"

	"github.com/RedisLabs/rediscloud-go-api/service/task"
)

type Client struct {
	Task *task.Api
}

func NewClient(configs ...Option) Client {
	config := &Options{
		baseUrl:   "https://api.redislabs.com/v1",
		userAgent: userAgent,
		apiKey:    os.Getenv("REDISLABS_API_KEY"),
		secretKey: os.Getenv("REDISLABS_SECRET_KEY"),
		transport: http.DefaultTransport,
	}

	for _, option := range configs {
		option(config)
	}

	client := &http.Client{
		Transport: config.roundTripper(),
	}
	t := task.NewApi(client, config.baseUrl)

	return Client{
		Task: t,
	}
}

type Options struct {
	baseUrl   string
	apiKey    string
	secretKey string
	userAgent string
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
