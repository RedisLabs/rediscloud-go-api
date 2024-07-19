package rediscloud_api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/redis_rules"
	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/roles"
	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/users"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/service/account"
	"github.com/RedisLabs/rediscloud-go-api/service/cloud_accounts"
	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	fixedDatabases "github.com/RedisLabs/rediscloud-go-api/service/fixed/databases"
	"github.com/RedisLabs/rediscloud-go-api/service/fixed/plans"
	"github.com/RedisLabs/rediscloud-go-api/service/fixed/plans/plan_subscriptions"
	fixedSubscriptions "github.com/RedisLabs/rediscloud-go-api/service/fixed/subscriptions"
	"github.com/RedisLabs/rediscloud-go-api/service/latest_backups"
	"github.com/RedisLabs/rediscloud-go-api/service/latest_imports"
	"github.com/RedisLabs/rediscloud-go-api/service/maintenance"
	"github.com/RedisLabs/rediscloud-go-api/service/pricing"
	"github.com/RedisLabs/rediscloud-go-api/service/regions"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
	"github.com/RedisLabs/rediscloud-go-api/service/transit_gateway/attachments"
)

type Client struct {
	Account                   *account.API
	CloudAccount              *cloud_accounts.API
	Database                  *databases.API
	Subscription              *subscriptions.API
	Regions                   *regions.API
	LatestBackup              *latest_backups.API
	LatestImport              *latest_imports.API
	Maintenance               *maintenance.API
	Pricing                   *pricing.API
	TransitGatewayAttachments *attachments.API
	// fixed
	FixedPlans             *plans.API
	FixedSubscriptions     *fixedSubscriptions.API
	FixedPlanSubscriptions *plan_subscriptions.API
	FixedDatabases         *fixedDatabases.API
	// acl
	RedisRules *redis_rules.API
	Roles      *roles.API
	Users      *users.API
}

func NewClient(configs ...Option) (*Client, error) {
	config := &Options{
		baseUrl:   "https://api.redislabs.com/v1",
		userAgent: userAgent,
		apiKey:    os.Getenv(AccessKeyEnvVar),
		secretKey: os.Getenv(SecretKeyEnvVar),
		logger:    &defaultLogger{},
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

	t := internal.NewAPI(client, config.logger)

	return &Client{
		Account:                   account.NewAPI(client),
		CloudAccount:              cloud_accounts.NewAPI(client, t, config.logger),
		Database:                  databases.NewAPI(client, t, config.logger),
		Subscription:              subscriptions.NewAPI(client, t, config.logger),
		Regions:                   regions.NewAPI(client, t, config.logger),
		LatestBackup:              latest_backups.NewAPI(client, t, config.logger),
		LatestImport:              latest_imports.NewAPI(client, t, config.logger),
		Maintenance:               maintenance.NewAPI(client, t, config.logger),
		Pricing:                   pricing.NewAPI(client),
		TransitGatewayAttachments: attachments.NewAPI(client, t, config.logger),
		// fixed
		FixedPlans:             plans.NewAPI(client, config.logger),
		FixedPlanSubscriptions: plan_subscriptions.NewAPI(client, config.logger),
		FixedSubscriptions:     fixedSubscriptions.NewAPI(client, t, config.logger),
		FixedDatabases:         fixedDatabases.NewAPI(client, t, config.logger),
		// acl
		RedisRules: redis_rules.NewAPI(client, t, config.logger),
		Roles:      roles.NewAPI(client, t, config.logger),
		Users:      users.NewAPI(client, t, config.logger),
	}, nil
}

type Options struct {
	baseUrl     string
	apiKey      string
	secretKey   string
	userAgent   string
	logger      Log
	transport   http.RoundTripper
	logRequests bool
}

func (o Options) roundTripper() http.RoundTripper {
	return &credentialTripper{
		apiKey:      o.apiKey,
		secretKey:   o.secretKey,
		wrapped:     o.transport,
		logRequests: o.logRequests,
		logger:      o.logger,
		userAgent:   o.userAgent,
	}
}

type Option func(*Options)

// Auth is used to set the authentication credentials - will otherwise default to using environment variables
// for the credentials.
func Auth(apiKey string, secretKey string) Option {
	return func(options *Options) {
		options.apiKey = apiKey
		options.secretKey = secretKey
	}
}

// BaseURL sets the URL to use for the API endpoint - will default to `https://api.redislabs.com/v1`.
func BaseURL(url string) Option {
	return func(options *Options) {
		options.baseUrl = url
	}
}

// LogRequests allows the logging of HTTP request and responses - will default to false (disabled).
func LogRequests(enable bool) Option {
	return func(options *Options) {
		options.logRequests = enable
	}
}

// Transporter allows the customisation of the RoundTripper used to communicate with the API - will default to the
// Go default.
func Transporter(transporter http.RoundTripper) Option {
	return func(options *Options) {
		options.transport = transporter
	}
}

// AdditionalUserAgent allows extra information to be appended to the user agent sent in all requests to the API.
func AdditionalUserAgent(additional string) Option {
	return func(options *Options) {
		options.userAgent += " " + additional
	}
}

// Logger allows for a custom implementation to handle the debug log messages - defaults to using the Go standard log
// package.
func Logger(log Log) Option {
	return func(options *Options) {
		options.logger = log
	}
}

type Log interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type defaultLogger struct{}

func (d *defaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (d *defaultLogger) Println(v ...interface{}) {
	log.Println(v...)
}

type credentialTripper struct {
	apiKey      string
	secretKey   string
	wrapped     http.RoundTripper
	logRequests bool
	logger      Log
	userAgent   string
}

func (c *credentialTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", c.userAgent)

	if c.logRequests {
		data, _ := httputil.DumpRequestOut(request, true)
		if data != nil {
			c.logger.Printf(`DEBUG: Request %s:
---[ REQUEST ]---
%s`, escapePath(request.URL.Path), redactPasswords(prettyPrint(data)))
		}
	}

	// Credentials added _after_ the request was logged to avoid accidentally logging them
	request.Header.Set("X-Api-Key", c.apiKey)
	request.Header.Set("X-Api-Secret-Key", c.secretKey)

	response, err := c.wrapped.RoundTrip(request)
	if err != nil {
		return response, err
	}

	if c.logRequests {
		data, _ := httputil.DumpResponse(response, true)
		if data != nil {
			c.logger.Printf(`DEBUG: Response %s:
---[ RESPONSE ]---
%s`, escapePath(request.URL.Path), redactPasswords(prettyPrint(data)))
		}
	}
	return response, nil
}

func prettyPrint(data []byte) string {
	lines := strings.Split(string(data), "\n")
	// A JSON body that wasn't indented would have ended up as a single line in the dumped information,
	// so try and find a line which is valid JSON and then indent it
	for i, line := range lines {
		asBytes := []byte(line)
		if json.Valid(asBytes) {
			var indented bytes.Buffer
			if err := json.Indent(&indented, asBytes, "", "  "); err == nil {
				lines[i] = indented.String()
			}
		}
	}
	return strings.Join(lines, "\n")
}

// redactPasswords: Redacts password values from a JSON message.
func redactPasswords(data string) string {
	m1 := regexp.MustCompile(`"password"\s*:\s*"(?:[^"\\]|\\.)*"`)
	output := m1.ReplaceAllString(data, "\"password\": \"REDACTED\"")
	m2 := regexp.MustCompile(`"global_password"\s*:\s*"(?:[^"\\]|\\.)*"`)
	return m2.ReplaceAllString(output, "\"global_password\": \"REDACTED\"")
}

func escapePath(path string) string {
	escapedPath := strings.Replace(path, "\n", "", -1)
	escapedPath = strings.Replace(escapedPath, "\r", "", -1)
	return escapedPath
}

var _ http.RoundTripper = &credentialTripper{}
