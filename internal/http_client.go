package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClient struct {
	client  *http.Client
	baseUrl *url.URL
}

func NewHttpClient(client *http.Client, baseUrl string) (*HttpClient, error) {
	parsed, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &HttpClient{client: client, baseUrl: parsed}, nil
}

func (c *HttpClient) Get(ctx context.Context, name, path string, responseBody interface{}) error {
	return c.withoutRequestBody(ctx, http.MethodGet, name, path, responseBody)
}

func (c *HttpClient) Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	return c.withRequestBody(ctx, http.MethodPut, name, path, requestBody, responseBody)
}

func (c *HttpClient) Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	return c.withRequestBody(ctx, http.MethodPost, name, path, requestBody, responseBody)
}

func (c *HttpClient) Delete(ctx context.Context, name, path string, responseBody interface{}) error {
	return c.withoutRequestBody(ctx, http.MethodDelete, name, path, responseBody)
}

func (c *HttpClient) withoutRequestBody(ctx context.Context, method, name, path string, responseBody interface{}) error {
	parsed := new(url.URL)
	*parsed = *c.baseUrl

	parsed.Path += path

	u := parsed.String()

	request, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return fmt.Errorf("failed to create request to %s: %w", name, err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to %s: %w", name, err)
	}

	defer response.Body.Close()

	if response.StatusCode > 299 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("failed to %s: %d - %s", name, response.StatusCode, body)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response to %s: %w", name, err)
	}

	return nil
}

func (c *HttpClient) withRequestBody(ctx context.Context, method, name, path string, requestBody interface{}, responseBody interface{}) error {
	parsed := new(url.URL)
	*parsed = *c.baseUrl

	parsed.Path += path

	u := parsed.String()

	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(requestBody); err != nil {
		return fmt.Errorf("failed to encode request for %s: %w", name, err)
	}

	request, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return fmt.Errorf("failed to create request to %s: %w", name, err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to %s: %w", name, err)
	}

	defer response.Body.Close()

	if response.StatusCode > 299 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("failed to %s: %d - %s", name, response.StatusCode, body)
	}

	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response to %s: %w", name, err)
	}

	return nil
}
