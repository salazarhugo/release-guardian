package app

import (
	"github.com/xanzy/go-gitlab"
	"os"
)

func NewGitlab() (*gitlab.Client, error) {
	git, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		return nil, err
	}

	return git, nil
}
