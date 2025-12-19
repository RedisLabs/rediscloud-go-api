package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/avast/retry-go/v4"
)

const (
	defaultWindowLimit    = 400
	defaultWindowDuration = 1 * time.Minute

	headerRateLimitRemaining = "X-Rate-Limit-Remaining"
)

type HttpClient struct {
	client           *http.Client
	baseUrl          *url.URL
	rateLimiter      RateLimiter
	retryEnabled     bool
	retryMaxDelay    time.Duration
	retryDelay       time.Duration
	retryMaxAttempts uint
	logger           Log
}

func NewHttpClient(client *http.Client, baseUrl string, logger Log) (*HttpClient, error) {
	parsed, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return &HttpClient{
		client:           client,
		baseUrl:          parsed,
		rateLimiter:      newFixedWindowCountRateLimiter(defaultWindowLimit, defaultWindowDuration),
		retryEnabled:     true,
		retryMaxAttempts: 10,
		retryDelay:       1 * time.Second,
		retryMaxDelay:    defaultWindowDuration,
		logger:           logger,
	}, nil
}

func (c *HttpClient) Get(ctx context.Context, name, path string, responseBody interface{}) error {
	return c.connectionWithRetries(ctx, http.MethodGet, name, path, nil, nil, responseBody)
}

func (c *HttpClient) GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error {
	return c.connectionWithRetries(ctx, http.MethodGet, name, path, query, nil, responseBody)
}

func (c *HttpClient) Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	return c.connectionWithRetries(ctx, http.MethodPut, name, path, nil, requestBody, responseBody)
}

func (c *HttpClient) Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	return c.connectionWithRetries(ctx, http.MethodPost, name, path, nil, requestBody, responseBody)
}

func (c *HttpClient) Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	return c.connectionWithRetries(ctx, http.MethodDelete, name, path, nil, requestBody, responseBody)
}

func (c *HttpClient) DeleteWithQuery(ctx context.Context, name, path string, query url.Values, requestBody interface{}, responseBody interface{}) error {
	return c.connectionWithRetries(ctx, http.MethodDelete, name, path, query, requestBody, responseBody)
}

func (c *HttpClient) connectionWithRetries(ctx context.Context, method, name, path string, query url.Values, requestBody interface{}, responseBody interface{}) error {
	return retry.Do(func() error {
		return c.connection(ctx, method, name, path, query, requestBody, responseBody)
	},
		retry.Attempts(c.retryMaxAttempts),
		retry.Delay(c.retryDelay),
		retry.MaxDelay(c.retryMaxDelay),
		retry.RetryIf(func(err error) bool {
			if !c.retryEnabled {
				return false
			}
			var target *HTTPError
			if errors.As(err, &target) && target.StatusCode == http.StatusTooManyRequests {
				c.logger.Println(fmt.Sprintf("status code 429 received, request will be retried"))
				return true
			}
			return false
		}),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
	)
}

func (c *HttpClient) connection(ctx context.Context, method, name, path string, query url.Values, requestBody interface{}, responseBody interface{}) error {
	if c.rateLimiter != nil {
		err := c.rateLimiter.Wait(ctx)
		if err != nil {
			return err
		}
	}

	parsed := new(url.URL)
	*parsed = *c.baseUrl

	parsed.Path += path
	if query != nil {
		parsed.RawQuery = query.Encode()
	}

	u := parsed.String()

	var body io.Reader
	if requestBody != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(requestBody); err != nil {
			return fmt.Errorf("failed to encode request for %s: %w", name, err)
		}
		body = buf
	}

	request, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return fmt.Errorf("failed to create request to %s: %w", name, err)
	}

	// The API expects this entry in the header in all requests.
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to %s: %w", name, err)
	}

	remainingLimit := response.Header.Get(headerRateLimitRemaining)
	if remainingLimit != "" {
		if limit, err := strconv.Atoi(remainingLimit); err == nil {
			err = c.rateLimiter.Update(limit)
			if err != nil {
				return err
			}
		}
	}

	defer response.Body.Close()

	if response.StatusCode > 299 {
		body, _ := io.ReadAll(response.Body)
		return &HTTPError{
			Name:       name,
			StatusCode: response.StatusCode,
			Body:       body,
		}
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response to %s: %w", name, err)
	}

	return nil
}

type HTTPError struct {
	Name       string
	StatusCode int
	Body       []byte
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("failed to %s: %d - %s", h.Name, h.StatusCode, h.Body)
}

var _ error = &HTTPError{}
