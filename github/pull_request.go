package github

import (
	"context"
	"strconv"

	"github.com/google/go-github/v64/github"
)

// CheckRunCompletion checks the completion status of check runs for a given pull request.
// It returns true if all check runs have completed successfully, false if any check run has failed,
// and an error if there was an issue fetching the check runs.
//
// Parameters:
//   - ctx: The context for the request.
//   - client: The GitHub client to use for making API requests.
//   - owner: The owner of the repository.
//   - repo: The name of the repository.
//   - prNumber: The pull request number.
//
// Returns:
//   - bool: True if all check runs are successful, false otherwise.
//   - error: An error if there was an issue fetching the check runs.
func CheckRunCompletion(ctx context.Context, client *github.Client, owner, repo string, prNumber int) (bool, error) {
	prNumberStr := strconv.Itoa(prNumber)

	checks, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repo, "refs/pull/"+prNumberStr+"/head", nil)
	if err != nil {
		return false, err
	}
	count := 0
	completed := false

	for !completed {
		for _, check := range checks.CheckRuns {
			if check.GetStatus() == "queued" || check.GetStatus() == "in_progress" {
				completed = false
				count++
			} else if check.GetStatus() == "completed" {
				completed = true
				if completed {
					if check.GetConclusion() == "success" {
						return true, nil
					} else if check.GetConclusion() == "failure" {
						return false, nil
					}
					break
				}
			}
		}
	}
	return false, nil
}
