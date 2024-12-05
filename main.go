package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pmclain/toggl-jira/internal/jira"
	"github.com/pmclain/toggl-jira/internal/toggl"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Validate environment variables
	requiredEnvVars := []string{"TOGGL_API_TOKEN", "JIRA_HOST", "JIRA_EMAIL", "JIRA_API_TOKEN"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
	}

	// Get time entries from the last 48 hours
	startTime := time.Now().AddDate(0, 0, -2)
	log.Printf("Starting sync from %s", startTime.Format(time.RFC3339))
	
	togglClient := toggl.NewClient(os.Getenv("TOGGL_API_TOKEN"))
	jiraClient := jira.NewClient(
		os.Getenv("JIRA_HOST"),
		os.Getenv("JIRA_EMAIL"),
		os.Getenv("JIRA_API_TOKEN"),
	)

	entries, err := togglClient.GetTimeEntries(startTime)
	if err != nil {
		log.Fatalf("Error getting time entries: %v", err)
	}

	if len(entries) == 0 {
		log.Println("No time entries found in the last 48 hours")
		return
	}

	log.Printf("Processing %d time entries", len(entries))
	for _, entry := range entries {
		log.Printf("Processing entry: %s", entry.Description)
		if err := jiraClient.AddWorkLog(entry); err != nil {
			log.Printf("Error adding work log for entry %s: %v", entry.Description, err)
			continue
		}
		log.Printf("Successfully added work log for entry: %s", entry.Description)
	}
} 