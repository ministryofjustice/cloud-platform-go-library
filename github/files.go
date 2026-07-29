package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v64/github"
)

// SelectFile checks if the given GitHub commit file is a Terraform file
// located in the "namespaces/live" directory. If both conditions are met,
// it returns the file; otherwise, it returns nil.
//
// Parameters:
//   - file: A pointer to a github.CommitFile object representing the file to be checked.
//   - A pointer to the github.CommitFile object if the file meets the criteria, or nil otherwise.e during decoding.
//
// Returns:
//   - A pointer to the github.CommitFile object if the file meets the criteria, or nil otherwise.
func SelectFile(file *github.CommitFile) *github.CommitFile {
	if strings.Contains(*file.Filename, "namespaces/live") && strings.Contains(*file.Filename, ".tf") {
		return file
	} else {
		return nil
	}
}

// GetFileContent retrieves the content of a specified file from a GitHub repository.
//
// Parameters:
//   - client: A GitHub client instance used to interact with the GitHub API.
//   - ctx: The context for the request, which can be used to control timeouts and cancellations.
//   - file: A pointer to a github.CommitFile object representing the file to retrieve.
//   - owner: The owner of the repository.
//   - repo: The name of the repository.
//   - ref: The name of the commit/branch/tag.
//
// Returns:
//   - A pointer to a github.RepositoryContent object containing the file content.
//   - An error if there is any issue in fetching the file content.		fmt.Printf("Error decoding file content: %v\n", err)
func GetFileContent(client *github.Client, ctx context.Context, file *github.CommitFile, owner, repo, ref string) (*github.RepositoryContent, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *file.Filename, opts)
	if err != nil {
		fmt.Printf("Error fetching file content: %v\n", err)
		return nil, err
	}

	return content, nil
}

// DecodeContent decodes the content of a given GitHub repository file.
// It takes a pointer to a github.RepositoryContent object and returns the decoded content as a string.
// If an error occurs during decoding, it returns an empty string and the error.
//
// Parameters:
//   - content: A pointer to a github.RepositoryContent object representing the file content to decode.
//
// Returns:
//   - A string containing the decoded content of the file.
//   - An error if there is an issue during decoding.
func DecodeContent(content *github.RepositoryContent) (string, error) {
	decodeContent, err := content.GetContent()
	if err != nil {
		fmt.Printf("Error decoding file content: %v\n", err)
		return "", err
	}

	return decodeContent, nil
}
