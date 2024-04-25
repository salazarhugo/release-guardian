package app

import (
	"github.com/andygrunwald/go-jira"
	"github.com/xanzy/go-gitlab"
	"log"
	"strings"
)

func Print(
	issues []jira.Issue,
	commits []*gitlab.Commit,
) error {
	log.Printf("Total issues: %d", len(issues))
	log.Println()

	testIssues := filter(issues, func(issue jira.Issue) bool {
		return strings.Contains(issue.Fields.Type.Name, "Test")
	})

	log.Printf("Test issues: %d", len(testIssues))
	PrintIssues(testIssues, false)

	mergedIssues := getMergedIssues(issues, commits)
	log.Printf("Merged issues: %d", len(mergedIssues))
	PrintIssues(mergedIssues, false)

	remainingIssues := getRemainingIssues(getRemainingIssues(issues, mergedIssues), testIssues)
	log.Printf("Pending issues: %d", len(remainingIssues))
	PrintIssues(remainingIssues, true)

	log.Printf("Found %d active tickets | %d are merged into develop", len(issues)-len(testIssues), len(mergedIssues))

	return nil
}

func PrintIssues(issues []jira.Issue, status bool) {
	for _, issue := range issues {
		if status {
			if issue.Fields.Status.Name == "In Progress" {
				log.Printf("	%s | %s | %s", issue.Key, issue.Fields.Status.Name, issue.Fields.Assignee.Name)
			} else {
				log.Printf("	%s | %s", issue.Key, issue.Fields.Status.Name)
			}
		} else {
			log.Printf("	%s", issue.Key)
		}
	}
	log.Println()
}

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func getMergedIssues(
	issues []jira.Issue,
	commits []*gitlab.Commit,
) []jira.Issue {
	res := make([]jira.Issue, 0)
	for _, issue := range issues {
		for _, commit := range commits {
			if strings.Contains(commit.Message, issue.Key) {
				res = append(res, issue)
				break
			}
		}
	}

	return res
}

func getRemainingIssues(
	issues []jira.Issue,
	mergedIssues []jira.Issue,
) []jira.Issue {
	var remainingIssues []jira.Issue
	for _, issue := range issues {
		merged := false
		for _, mergedIssue := range mergedIssues {
			if issue.ID == mergedIssue.ID {
				merged = true
				break
			}
		}
		if !merged {
			remainingIssues = append(remainingIssues, issue)
		}
	}
	return remainingIssues
}
