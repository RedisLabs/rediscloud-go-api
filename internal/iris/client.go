// Package iris provides the HTTP transport shared by the Iris service connectors.
package iris

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient is implemented by http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	baseURL *url.URL
	token   string
	client  HTTPClient
}

// NewClient expects a service origin, optionally including a reverse-proxy prefix.
// Operation paths include their own API version.
func NewClient(baseURL, token string, client HTTPClient) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Iris base URL: %w", err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return nil, errors.New("iris base URL must be an HTTP(S) origin with an optional path prefix")
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{baseURL: u, token: token, client: client}, nil
}

// HTTPError retains the status, headers and problem response for errors.As callers.
// Error omits response contents, which may contain sensitive information.
type HTTPError struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("Iris API returned HTTP %d", e.StatusCode)
}

// Call performs one synchronous operation. Mutations are never retried implicitly.
func (c *Client) Call(ctx context.Context, method, path string, query url.Values, headers http.Header, body, result interface{}) (http.Header, error) {
	u := *c.baseURL
	escaped := strings.TrimRight(u.EscapedPath(), "/") + path
	decoded, err := url.PathUnescape(escaped)
	if err != nil {
		return nil, fmt.Errorf("invalid Iris operation path: %w", err)
	}
	u.Path, u.RawPath, u.RawQuery = decoded, escaped, query.Encode()
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode Iris request: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("create Iris request: %w", err)
	}
	for name, values := range headers {
		req.Header[name] = append([]string(nil), values...)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	response, err := c.client.Do(req) //nolint:gosec // G704: origin is explicitly configured by the SDK caller
	if err != nil {
		return nil, fmt.Errorf("call Iris API: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			return response.Header, fmt.Errorf("read Iris error response: %w", err)
		}
		return response.Header, &HTTPError{StatusCode: response.StatusCode, Header: response.Header.Clone(), Body: data}
	}
	if result != nil {
		if err := json.NewDecoder(response.Body).Decode(result); err != nil {
			return response.Header, fmt.Errorf("decode Iris response: %w", err)
		}
	}
	return response.Header, nil
}

// Requester allows service connectors to share a transport or use a test double.
type Requester interface {
	Call(context.Context, string, string, url.Values, http.Header, interface{}, interface{}) (http.Header, error)
}

// Result decodes a JSON response and leaves error handling to the caller.
func Result[T any](ctx context.Context, client Requester, method, path string, query url.Values, headers http.Header, body interface{}) (*T, http.Header, error) {
	var result T
	h, err := client.Call(ctx, method, path, query, headers, body, &result)
	if err != nil {
		return nil, h, err
	}
	return &result, h, nil
}
