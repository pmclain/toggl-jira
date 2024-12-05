package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pmclain/toggl-jira/internal/toggl"
)

func TestExtractIssuesFromTimeEntry(t *testing.T) {
	tests := []struct {
		name        string
		description string
		projects    string
		want        []string
	}{
		{
			name:        "single issue no filter",
			description: "ISSUE-52 doing work",
			projects:    "",
			want:        []string{"ISSUE-52"},
		},
		{
			name:        "single issue with matching filter",
			description: "ISSUE-52 doing work",
			projects:    "ISSUE",
			want:        []string{"ISSUE-52"},
		},
		{
			name:        "single issue with non-matching filter",
			description: "ISSUE-52 doing work",
			projects:    "OTHER",
			want:        nil,
		},
		{
			name:        "multiple issues",
			description: "ISSUE-52,ISSUE-55 doing work",
			projects:    "ISSUE",
			want:        []string{"ISSUE-52", "ISSUE-55"},
		},
		{
			name:        "no issue",
			description: "doing work",
			projects:    "ISSUE",
			want:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("JIRA_PROJECTS", tt.projects)
			client := NewClient("host", "email", "token")
			got := client.extractIssuesFromTimeEntry(tt.description)
			if len(got) != len(tt.want) {
				t.Errorf("extractIssuesFromTimeEntry() got %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractIssuesFromTimeEntry() got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAddWorkLog(t *testing.T) {
	tests := []struct {
		name           string
		entry         toggl.TimeEntry
		existingWorkLog string
		wantStatus    int
		wantError     bool
	}{
		{
			name: "create new worklog",
			entry: toggl.TimeEntry{
				ID:          1,
				Description: "ISSUE-52 doing work",
				Duration:    3600,
				Start:       time.Now(),
				Stop:        time.Now().Add(time.Hour),
			},
			existingWorkLog: `{"worklogs": []}`,
			wantStatus:      http.StatusCreated,
			wantError:       false,
		},
		{
			name: "skip running timer",
			entry: toggl.TimeEntry{
				ID:          1,
				Description: "ISSUE-52 doing work",
				Duration:    -1,
			},
			wantError: false,
		},
		{
			name: "no issue key",
			entry: toggl.TimeEntry{
				ID:          1,
				Description: "doing work",
				Duration:    3600,
				Start:       time.Now(),
				Stop:        time.Now().Add(time.Hour),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			
			// Mock the worklog endpoint
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				// Check auth
				username, password, ok := r.BasicAuth()
				if !ok || username != "email" || password != "token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				switch r.Method {
				case "GET":
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.existingWorkLog))
				case "POST":
					w.WriteHeader(tt.wantStatus)
				}
			})

			server := httptest.NewServer(mux)
			defer server.Close()

			client := NewClient(strings.TrimPrefix(server.URL, "http://"), "email", "token")
			client.SetScheme("http") // Use HTTP for testing
			err := client.AddWorkLog(tt.entry)

			if (err != nil) != tt.wantError {
				t.Errorf("AddWorkLog() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestGetWorkLogsForIssue(t *testing.T) {
	workLogs := JiraWorkLogResponse{
		Worklogs: []JiraWorkLog{
			{
				ID:      "1",
				Comment: "TogglID: 123 doing work",
			},
			{
				ID:      "2",
				Comment: "TogglID: 456 other work",
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check auth
		username, password, ok := r.BasicAuth()
		if !ok || username != "email" || password != "token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workLogs)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(strings.TrimPrefix(server.URL, "http://"), "email", "token")
	client.SetScheme("http") // Use HTTP for testing

	tests := []struct {
		name     string
		togglID  int64
		wantID   string
		wantNil  bool
	}{
		{
			name:    "existing worklog",
			togglID: 123,
			wantID:  "1",
			wantNil: false,
		},
		{
			name:    "no matching worklog",
			togglID: 789,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.getWorkLogsForIssue("ISSUE-1", tt.togglID)
			if err != nil {
				t.Errorf("getWorkLogsForIssue() error = %v", err)
				return
			}

			if tt.wantNil && got != nil {
				t.Errorf("getWorkLogsForIssue() got = %v, want nil", got)
			} else if !tt.wantNil && (got == nil || got.ID != tt.wantID) {
				t.Errorf("getWorkLogsForIssue() got = %v, want ID %s", got, tt.wantID)
			}
		})
	}
} 