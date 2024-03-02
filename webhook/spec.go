package webhook

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/github"
)

type FileContent struct {
	Path   string
	Change string
}

func getFileContents(ctx context.Context, client *github.Client, owner, repo string, paths []string) ([]FileContent, error) {
	var contents []FileContent
	for _, path := range paths {
		fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get content for path %s: %w", path, err)
		}

		content, err := base64.StdEncoding.DecodeString(*fileContent.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to decode file content %s: %w", path, err)
		}

		contents = append(contents, FileContent{
			path,
			string(content),
		})

	}
	return contents, nil

}
