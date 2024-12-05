package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pmclain/toggl-jira/internal/toggl"
)

var (
	issueKeyPattern = regexp.MustCompile(`^([A-Z]+-\d+)`)
)

type Client struct {
	host           string
	email          string
	apiToken       string
	supportedKeys  []string
	http           *http.Client
	scheme         string
}

type WorkLog struct {
	ID              string `json:"id,omitempty"`
	Comment         string `json:"comment"`
	TimeSpent       string `json:"timeSpent"`
	Started         string `json:"started"`
	TimeSpentSeconds int    `json:"timeSpentSeconds,omitempty"`
}

func NewClient(host, email, apiToken string) *Client {
	// Remove any protocol prefix from host
	host = strings.TrimPrefix(strings.TrimPrefix(host, "https://"), "http://")
	
	// Get supported project keys
	supportedKeys := getSupportedJiraKeys()
	if len(supportedKeys) > 0 {
		log.Printf("Supported JIRA projects: %s", strings.Join(supportedKeys, ", "))
	} else {
		log.Println("No JIRA project filter configured, accepting all projects")
	}

	return &Client{
		host:          host,
		email:         email,
		apiToken:      apiToken,
		supportedKeys: supportedKeys,
		http:         &http.Client{Timeout: 10 * time.Second},
		scheme:       "https", // Default to HTTPS
	}
}

// SetScheme allows overriding the URL scheme for testing
func (c *Client) SetScheme(scheme string) {
	c.scheme = scheme
}

func (c *Client) buildURL(path string) string {
	return fmt.Sprintf("%s://%s%s", c.scheme, c.host, path)
}

func (c *Client) AddWorkLog(entry toggl.TimeEntry) error {
	if entry.Duration < 1 || entry.Stop.IsZero() {
		log.Printf("Skipping running timer: %s", entry.Description)
		return nil
	}

	issueKeys := c.extractIssuesFromTimeEntry(entry.Description)
	if len(issueKeys) == 0 {
		log.Printf("No supported issue found in description: %s", entry.Description)
		return nil
	}

	// For now, we'll use the first issue key found (TODO: support multiple issues)
	issueKey := issueKeys[0]

	// Convert duration from seconds to minutes
	minutes := int(math.Ceil(float64(entry.Duration) / 60.0))

	// Format the time in Jira's expected format with UTC timezone
	started := entry.Start.UTC().Format("2006-01-02T15:04:05.000+0000")

	workLog := WorkLog{
		TimeSpent: fmt.Sprintf("%dm", minutes),
		Started:   started,
		Comment:   fmt.Sprintf("TogglID: %d %s", entry.ID, strings.TrimSpace(strings.TrimPrefix(entry.Description, issueKey))),
	}

	// Check for existing worklog
	existingWorkLog, err := c.getWorkLogsForIssue(issueKey, entry.ID)
	if err != nil {
		return fmt.Errorf("checking existing worklog: %w", err)
	}

	if existingWorkLog != nil {
		return c.updateWorkLog(issueKey, existingWorkLog.ID, workLog)
	}

	return c.createWorkLog(issueKey, workLog)
}

func (c *Client) createWorkLog(issueKey string, workLog WorkLog) error {
	payload, err := json.Marshal(workLog)
	if err != nil {
		return fmt.Errorf("marshaling work log: %w", err)
	}

	url := c.buildURL(fmt.Sprintf("/rest/api/latest/issue/%s/worklog?notifyUsers=false", issueKey))
	log.Printf("Creating worklog for issue %s", issueKey)
	log.Printf("Payload: %s", string(payload))
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully created worklog for issue %s", issueKey)
	return nil
}

func (c *Client) updateWorkLog(issueKey, worklogID string, workLog WorkLog) error {
	payload, err := json.Marshal(workLog)
	if err != nil {
		return fmt.Errorf("marshaling work log: %w", err)
	}

	url := c.buildURL(fmt.Sprintf("/rest/api/latest/issue/%s/worklog/%s?notifyUsers=false", issueKey, worklogID))
	log.Printf("Updating worklog for issue %s", issueKey)
	log.Printf("Payload: %s", string(payload))

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Successfully updated worklog for issue %s", issueKey)
	return nil
}

type JiraWorkLog struct {
	ID              string `json:"id"`
	Comment         string `json:"comment"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
}

type JiraWorkLogResponse struct {
	Worklogs []JiraWorkLog `json:"worklogs"`
}

func (c *Client) getWorkLogsForIssue(issueKey string, togglID int64) (*JiraWorkLog, error) {
	url := c.buildURL(fmt.Sprintf("/rest/api/latest/issue/%s/worklog", issueKey))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.email, c.apiToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var response JiraWorkLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	togglIDStr := fmt.Sprintf("TogglID: %d", togglID)
	for _, workLog := range response.Worklogs {
		if strings.Contains(workLog.Comment, togglIDStr) {
			return &workLog, nil
		}
	}

	return nil, nil
}

func (c *Client) extractIssuesFromTimeEntry(description string) []string {
	if len(c.supportedKeys) == 0 {
		// If no supported keys are configured, use the default pattern
		matches := issueKeyPattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			return []string{matches[1]}
		}
		return nil
	}

	// Create a pattern that matches any of the supported project keys
	pattern := fmt.Sprintf(`(%s-\d+)`, strings.Join(c.supportedKeys, "|"))
	keyRegex := regexp.MustCompile(pattern)
	
	matches := keyRegex.FindAllString(description, -1)
	return matches
}

func getSupportedJiraKeys() []string {
	keys := strings.Split(os.Getenv("JIRA_PROJECTS"), ",")
	var result []string
	for _, key := range keys {
		if trimmed := strings.TrimSpace(key); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
} 