package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"release-guardian/internal/app"
)

func searchJIRA(
	client *jira.Client,
	jql string,
) ([]jira.Issue, error) {
	// Define the Jira search options
	opts := &jira.SearchOptions{
		StartAt:       0,
		MaxResults:    500,
		Expand:        "",
		Fields:        []string{"key", "resolution", "issuetype", "assignee", "status"},
		ValidateQuery: "",
	}

	// Perform the JQL search
	issues, _, err := client.Issue.Search(jql, opts)

	// Handle potential errors
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func init() {
	log.SetFlags(0)
}

func main() {
	// Replace with your JIRA server base URL and credentials
	jiraServer := "https://jira.agilefabric.fr.carrefour.com"

	tp := jira.BearerAuthTransport{
		Token:     os.Getenv("JIRA_TOKEN"),
		Transport: nil,
	}

	// Create a new JIRA client
	client, err := jira.NewClient(tp.Client(), jiraServer)
	if err != nil {
		fmt.Println("Error creating JIRA client:", err)
		return
	}

	// Define your JQL search string (e.g. project = MYPROJECT AND fixVersion = 1.0)
	jql := "project = DFEM AND fixVersion = 53263"

	// Call the searchJIRA function
	issues, err := searchJIRA(client, jql)
	if err != nil {
		fmt.Println("Error searching JIRA:", err)
		return
	}

	git, err := app.NewGitlab()
	if err != nil {
		return
	}
	develop := "develop"

	commits, _, err := git.Commits.ListCommits(49153483, &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: 100,
		},
		RefName: &develop,
	})
	if err != nil {
		return
	}

	err = app.Print(issues, commits)
	if err != nil {
		return
	}
}
