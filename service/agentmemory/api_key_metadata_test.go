package agentmemory

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestAPI_RejectsKeyMetadataWithoutID(t *testing.T) {
	for _, operation := range []string{"list", "get", "create-before-post", "create-after-post"} {
		for _, layout := range []string{"missing-only", "missing-before-match", "missing-after-match"} {
			t.Run(operation+"/"+layout, func(t *testing.T) {
				listCalls, postCalls := 0, 0
				metadata := APIKeyInfo{KeyName: "runtime", ObfuscatedToken: "FAKE-METADATA"}
				matching := APIKeyInfo{APIKeyID: "key-1", KeyName: "runtime"}
				bad := []APIKeyInfo{metadata}
				switch layout {
				case "missing-before-match":
					bad = append(bad, matching)
				case "missing-after-match":
					bad = []APIKeyInfo{matching, metadata}
				}
				client := &fakeHTTPClient{
					getFunc: func(_ context.Context, _, _ string, responseBody interface{}) error {
						listCalls++
						response := responseBody.(*ListAPIKeysResponse)
						response.APIKeys = bad
						if operation == "create-after-post" && listCalls == 1 {
							response.APIKeys = []APIKeyInfo{}
						}
						return nil
					},
					postFunc: func(_ context.Context, _, _ string, _, responseBody interface{}) error {
						postCalls++
						responseBody.(*CreateAPIKeyResponse).APIKey = "FAKE-CREATED-SECRET"
						return nil
					},
				}
				subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})
				var err error
				switch operation {
				case "list":
					var keys []APIKeyInfo
					keys, err = subject.ListAPIKeys(context.Background(), "store-1")
					if keys != nil {
						t.Error("malformed list returned usable metadata")
					}
				case "get":
					var key *APIKeyInfo
					key, err = subject.GetAPIKey(context.Background(), "store-1", "key-1")
					if key != nil {
						t.Error("malformed list returned a key")
					}
				default:
					var key *CreatedAPIKey
					key, err = subject.CreateAPIKey(context.Background(), "store-1", "runtime")
					if key != nil {
						t.Error("malformed list returned a created key with unreliable identity")
					}
				}
				if err == nil {
					t.Fatal("metadata without an ID was accepted")
				}
				var notFound *APIKeyNotFound
				if errors.As(err, &notFound) {
					t.Error("malformed metadata was interpreted as confirmed absence")
				}
				if !strings.Contains(err.Error(), "metadata without an ID") {
					t.Errorf("expected metadata validation error, got %v", err)
				}
				if strings.Contains(err.Error(), "FAKE-") {
					t.Error("metadata validation error included response contents")
				}
				wantLists, wantPosts := 1, 0
				if operation == "create-after-post" {
					wantLists, wantPosts = 2, 1
				}
				if listCalls != wantLists || postCalls != wantPosts {
					t.Errorf("unexpected API calls: list=%d post=%d", listCalls, postCalls)
				}
			})
		}
	}
}

func TestAPI_ValidKeyMetadataControls(t *testing.T) {
	for _, scenario := range []string{"empty-list", "other-key", "matching-key", "id-without-name"} {
		t.Run(scenario, func(t *testing.T) {
			keys := []APIKeyInfo{}
			switch scenario {
			case "other-key":
				keys = []APIKeyInfo{{APIKeyID: "key-2", KeyName: "other"}}
			case "matching-key":
				keys = []APIKeyInfo{{APIKeyID: "key-1", KeyName: "runtime"}}
			case "id-without-name":
				// Name validation belongs to consumers; this check only establishes identity.
				keys = []APIKeyInfo{{APIKeyID: "key-1"}}
			}
			client := &fakeHTTPClient{getFunc: func(_ context.Context, _, _ string, responseBody interface{}) error {
				responseBody.(*ListAPIKeysResponse).APIKeys = keys
				return nil
			}}
			subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})
			listed, err := subject.ListAPIKeys(context.Background(), "store-1")
			if err != nil || len(listed) != len(keys) {
				t.Fatal("valid IDs were rejected by ListAPIKeys")
			}
			key, err := subject.GetAPIKey(context.Background(), "store-1", "key-1")
			if scenario == "empty-list" || scenario == "other-key" {
				var notFound *APIKeyNotFound
				if !errors.As(err, &notFound) || key != nil {
					t.Error("valid absence no longer returns APIKeyNotFound")
				}
			} else if err != nil || key == nil || key.APIKeyID != "key-1" {
				t.Error("valid matching key was not returned")
			}
		})
	}
}
