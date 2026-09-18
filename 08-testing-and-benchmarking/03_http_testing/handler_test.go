package httptesting

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreetHandler(t *testing.T) {
	tests := []struct {
		name         string
		queryName    string
		wantGreeting string
	}{
		{"With explicit name", "Divya", "Hello, Divya!"},
		{"Without name parameter (fallback)", "", "Hello, Guest!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/greet"
			if tt.queryName != "" {
				url += "?name=" + tt.queryName
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			GreetHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", resp.StatusCode)
			}

			var body GreetResponse
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if body.Greeting != tt.wantGreeting {
				t.Errorf("expected greeting %q, got %q", tt.wantGreeting, body.Greeting)
			}
		})
	}
}
