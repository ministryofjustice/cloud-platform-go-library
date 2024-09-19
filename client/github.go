package client

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v64/github"
)

func GitHubClient(token string, ctx context.Context) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)

	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil
	}

	// Rate.Limit should most likely be 5000 when authorized.
	log.Printf("Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	return client
}
