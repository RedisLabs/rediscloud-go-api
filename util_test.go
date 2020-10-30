package rediscloud_api

import "net/http"

func testServer(path, apiKey, secretKey string, responses ...string) http.HandlerFunc {
	responseCount := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(500)
			return
		}
		if r.URL.Path != path {
			w.WriteHeader(501)
			return
		}
		if r.Header.Get("X-Api-Key") != apiKey {
			w.WriteHeader(502)
			return
		}
		if r.Header.Get("X-Api-Secret-Key") != secretKey {
			w.WriteHeader(503)
			return
		}

		body := responses[responseCount]
		responseCount++
		if responseCount > len(responses)-1 {
			responseCount = len(responses) - 2
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(body))
	}
}
