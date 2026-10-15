package rediscloud_api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientRejectsCrossOriginRedirectBeforeSendingCredentials(t *testing.T) {
	var contacted, credentialRequests atomic.Int32
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contacted.Add(1)
		if r.Header.Get("x-api-key") != "" || r.Header.Get("x-api-secret-key") != "" {
			credentialRequests.Add(1)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"storeId":"store-1","databaseId":"123","name":"test"}`))
	}))
	defer other.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/capture", http.StatusTemporaryRedirect)
	}))
	defer origin.Close()
	client, err := NewClient(BaseURL(origin.URL), Auth("fake-redirect-key", "fake-redirect-secret"))
	require.NoError(t, err)
	_, err = client.AgentMemory.Get(context.Background(), "store-1")
	assert.Error(t, err, "cross-origin redirect must be rejected")
	assert.Zero(t, contacted.Load(), "redirect target was contacted")
	assert.Zero(t, credentialRequests.Load(), "API credentials reached the redirect target")
}

func TestClientAllowsSameOriginRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/memory-stores/store-1" {
			http.Redirect(w, r, "/canonical/store-1", http.StatusTemporaryRedirect)
			return
		}
		assert.Equal(t, "fake-key", r.Header.Get("x-api-key"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"storeId":"store-1","databaseId":"123","name":"test"}`))
	}))
	defer server.Close()
	client, err := NewClient(BaseURL(server.URL), Auth("fake-key", "fake-secret"))
	require.NoError(t, err)
	store, err := client.AgentMemory.Get(context.Background(), "store-1")
	require.NoError(t, err)
	require.Equal(t, "store-1", store.StoreID)
}

func TestAPIClientRedirectOriginPolicy(t *testing.T) {
	original, err := url.Parse("https://api.example.test/resource")
	require.NoError(t, err)
	for _, test := range []struct {
		target string
		allow  bool
	}{
		{"https://API.example.test/next", true},
		{"https://api.example.test:443/next", true},
		{"https://api.example.test:444/next", false},
		{"http://api.example.test/next", false},
		{"https://other.example.test/next", false},
		{"https://user:password@api.example.test/next", false},
	} {
		t.Run(test.target, func(t *testing.T) {
			target, err := url.Parse(test.target)
			require.NoError(t, err)
			err = checkAPIClientRedirect(&http.Request{URL: target}, []*http.Request{{URL: original}})
			if test.allow {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
	via := make([]*http.Request, 10)
	for i := range via {
		via[i] = &http.Request{URL: original}
	}
	require.ErrorContains(t, checkAPIClientRedirect(&http.Request{URL: original}, via), "10")
}
