package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v64/github"
)

// GetPullRequestBranch retrieves the branch name associated with a given pull request.
//
// Parameters:
// - client: A GitHub client instance.
// - ctx: The context for the request.
// - owner: The owner of the repository.
// - repo: The name of the repository.
// - prNumber: The number of the pull request.
//
// Returns:
// - A string representing the branch name of the pull request.
// - An error if there is an issue fetching the pull request.
func GetPullRequestBranch(client *github.Client, ctx context.Context, owner, repo string, prNumber int) (string, error) {
	pull, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return "", fmt.Errorf("error fetching pull request: %w", err)
	}
	return *pull.Head.Ref, nil
}
