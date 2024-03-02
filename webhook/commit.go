package webhook

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/go-github/github"
)

func getMasterBranchSHA(ctx context.Context, client *github.Client, owner, repo string) (string, error) {
	ref, _, err := client.Git.GetRef(ctx, owner, repo, fmt.Sprintf("refs/heads/%s", targetBranchName))
	if err != nil {
		return "", fmt.Errorf("error getting master branch reference: %w", err)
	}
	return *ref.Object.SHA, nil
}

func getLatestCommitSHA(ctx context.Context, client *github.Client, owner, repo, branch string) (string, error) {
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err == nil {
		return *ref.Object.SHA, nil
	}

	if _, ok := err.(*github.ErrorResponse); ok {
		masterSHA, err := getMasterBranchSHA(ctx, client, owner, repo)
		if err != nil {
			return "", fmt.Errorf("error creating branch: %w", err)
		}
		// Branch doesn't exist, create a new branch
		ref, _, err = client.Git.CreateRef(ctx, owner, repo, &github.Reference{
			Ref: github.String("refs/heads/" + branch),
			Object: &github.GitObject{
				SHA: github.String(masterSHA), // SHA of the commit to base the branch off of
			},
		})
		if err != nil {
			return "", fmt.Errorf("error creating branch: %w", err)
		}
		return *ref.Object.SHA, nil
	}

	return "", fmt.Errorf("error getting reference: %w", err)
}

func CommitChange(change *ChangeInfo) error {

	ctx, client := GetGitHubClient()
	pullRequestOwner := change.Owner
	pullRequestRepo := change.Repo
	pullRequestChanges := change.Changes

	// Define the owner of the repository and the repository name
	owner := "abdallahy271"
	repo := "go-webhook"

	// Get the content of the file
	fileContents, err := getFileContents(ctx, client, pullRequestOwner, pullRequestRepo, pullRequestChanges)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Get the SHA of the latest commit on the branch
	latestCommitSHA, err := getLatestCommitSHA(ctx, client, owner, repo, change.SourceBranch)
	if err != nil {
		fmt.Println("Error getting latest commit:", err)
		return err
	}

	// Create a new tree with the changes you want to commit
	var entries []github.TreeEntry
	for _, fileContent := range fileContents {
		pullRequestOwnerPath := fmt.Sprintf("/%s", pullRequestOwner)
		prefixedPath := filepath.Join(pullRequestOwnerPath, fileContent.Path)

		entry := github.TreeEntry{
			Path:    github.String(prefixedPath),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String(fileContent.Change),
		}
		entries = append(entries, entry)
	}

	tree, _, err := client.Git.CreateTree(ctx, owner, repo, latestCommitSHA, entries)
	if err != nil {
		fmt.Println("Error creating tree:", err)
		return err
	}

	// Create a new commit using the new tree and the latest commit SHA
	newCommit, _, err := client.Git.CreateCommit(ctx, owner, repo, &github.Commit{
		Message: github.String("Commit message"),
		Tree:    tree,
		Parents: []github.Commit{{SHA: &latestCommitSHA}},
	})
	if err != nil {
		fmt.Println("Error creating commit:", err)
		return err
	}

	// Update the master branch reference to point to the new commit SHA
	_, _, err = client.Git.UpdateRef(ctx, owner, repo, &github.Reference{
		Ref: github.String("refs/heads/" + change.SourceBranch),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	}, false)
	if err != nil {
		fmt.Println("Error updating master branch reference:", err)
		return err
	}

	fmt.Printf("Commit to %s branch created successfully!", change.SourceBranch)
	return nil
}
