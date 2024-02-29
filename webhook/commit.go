package webhook

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func CommitChange() error {

	// Create a GitHub client with OAuth2 token authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Define the owner of the repository and the repository name
	owner := "abdallahy271"
	repo := "go-webhook"
	branch := "webhook"

	// Get the SHA of the latest commit on the branch
	ref, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		fmt.Println("Error getting reference:", err)
		return err
	}
	latestCommitSHA := *ref.Object.SHA

	// Create a new tree with the changes you want to commit
	// (This example just creates a dummy file)
	entries := []github.TreeEntry{
		{
			Path:    github.String("example2.txt"),
			Mode:    github.String("100644"),
			Type:    github.String("blob"),
			Content: github.String("This is another example file."),
		},
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
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA: newCommit.SHA,
		},
	}, false)
	if err != nil {
		fmt.Println("Error updating master branch reference:", err)
		return err
	}

	fmt.Println("Commit to master branch created successfully!")
	return nil
}
