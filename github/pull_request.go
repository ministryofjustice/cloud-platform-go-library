package github

import (
	"context"
	"strconv"

	"github.com/google/go-github/v64/github"
)

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
