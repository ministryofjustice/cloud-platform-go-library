package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v64/github"
)

func GetPullRequestBranch(client *github.Client, ctx context.Context, owner, repo string, prNumber int) (string, error) {
	pull, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return "", fmt.Errorf("error fetching pull request: %w", err)
	}
	return *pull.Head.Ref, nil
}
