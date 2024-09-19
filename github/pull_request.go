package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v64/github"
)

func GetPullRequestFiles(client *github.Client, ctx context.Context, o, r string, n int) ([]*github.CommitFile, *github.Response, error) {
	files, resp, err := client.PullRequests.ListFiles(ctx, o, r, n, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching files: %w", err)
	}

	return files, resp, err
}

func CheckRunCompletion(ctx context.Context, client *github.Client, owner, repo string, prNumber int) (bool, error) {
	// check status on CheckRun of a pull request
	// convert prNumber to string
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
