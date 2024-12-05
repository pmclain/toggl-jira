package toggl

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTimeEntries(t *testing.T) {
	mockEntries := []TimeEntry{
		{
			ID:          1951596187,
			Description: "ISSUE-52 doing work",
			Start:       time.Now().Add(-time.Hour),
			Stop:        time.Now(),
			Duration:    3600,
			WorkspaceID: 1391549,
		},
		{
			ID:          1951664141,
			Description: "ISSUE-112 daily",
			Start:       time.Now().Add(-2 * time.Hour),
			Stop:        time.Now().Add(-time.Hour),
			Duration:    3600,
			WorkspaceID: 1391549,
		},
	}

	tests := []struct {
		name       string
		status     int
		response   interface{}
		wantError  bool
		wantLength int
	}{
		{
			name:       "successful response",
			status:     http.StatusOK,
			response:   mockEntries,
			wantError:  false,
			wantLength: 2,
		},
		{
			name:       "empty response",
			status:     http.StatusOK,
			response:   []TimeEntry{},
			wantError:  false,
			wantLength: 0,
		},
		{
			name:      "unauthorized",
			status:    http.StatusUnauthorized,
			response:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Verify auth header
				_, password, ok := r.BasicAuth()
				if !ok || password != "api_token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// Verify query parameters
				if r.URL.Query().Get("start_date") == "" || r.URL.Query().Get("end_date") == "" {
					t.Errorf("Missing required query parameters")
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			client := NewClient("test-token")
			client.http = server.Client()
			// Override the API base URL for testing
			origAPIBase := togglAPIBase
			togglAPIBase = server.URL
			defer func() { togglAPIBase = origAPIBase }()

			entries, err := client.GetTimeEntries(time.Now().Add(-48 * time.Hour))

			if (err != nil) != tt.wantError {
				t.Errorf("GetTimeEntries() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && len(entries) != tt.wantLength {
				t.Errorf("GetTimeEntries() got %d entries, want %d", len(entries), tt.wantLength)
			}
		})
	}
} 