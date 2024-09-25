package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v64/github"
)

func SelectFile(file *github.CommitFile) *github.CommitFile {
	if strings.Contains(*file.Filename, "namespaces/live") && strings.Contains(*file.Filename, ".tf") {
		return file
	} else {
		return nil
	}
}

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

func DecodeContent(content *github.RepositoryContent) (string, error) {
	decodeContent, err := content.GetContent()
	if err != nil {
		fmt.Printf("Error decoding file content: %v\n", err)
		return "", err
	}

	return decodeContent, nil
}
