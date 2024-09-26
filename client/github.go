package client

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v64/github"
)

// GitHubClient is used to pass a GitHub token and context to methods that need to interact with the GitHub API.
// It returns a GitHub client that can be used to interact with the GitHub API.
func GitHubClient(token string, ctx context.Context) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)
	// This is a simple way to check if the token is valid.
	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil
	}

	// Rate.Limit should most likely be 5000 when authorized.
	// This is the number of requests you can make per hour.
	log.Printf("Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	// This is useful for tokens that expire.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	return client
}
