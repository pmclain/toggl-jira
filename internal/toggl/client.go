package toggl

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	togglAPIBase = "https://api.track.toggl.com/api/v9"
)

type Client struct {
	apiToken string
	http     *http.Client
}

type TimeEntry struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Duration    int       `json:"duration"`
	WorkspaceID int       `json:"workspace_id"`
	ProjectID   int       `json:"project_id,omitempty"`
}

func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetTimeEntries(startTime time.Time) ([]TimeEntry, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/me/time_entries", togglAPIBase), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.apiToken, "api_token")
	
	q := req.URL.Query()
	q.Add("start_date", startTime.Format(time.RFC3339))
	q.Add("end_date", time.Now().Format(time.RFC3339))
	req.URL.RawQuery = q.Encode()

	log.Printf("Fetching time entries from: %s", startTime.Format(time.RFC3339))
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Toggl Response: %s", string(body))

	var entries []TimeEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("decoding response: %w, body: %s", err, string(body))
	}

	log.Printf("Found %d time entries", len(entries))
	for _, entry := range entries {
		log.Printf("Entry: ID=%d, Description=%s, Duration=%d, Start=%s", 
			entry.ID, entry.Description, entry.Duration, entry.Start.Format(time.RFC3339))
	}

	return entries, nil
} 