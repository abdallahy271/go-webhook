package webhook

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/github"
)

func getFileContent(ctx context.Context, client *github.Client, owner, repo, path string) ([]byte, error) {
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get file content: %w", err)
	}

	content, err := base64.StdEncoding.DecodeString(*fileContent.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}
