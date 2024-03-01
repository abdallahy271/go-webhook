package webhook

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetGitHubClient() (context.Context, *github.Client) {
	// Create a GitHub client with OAuth2 token authentication
	ctx := context.Background()
	accessToken, _ := os.LookupEnv("GH_ACCESS_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return ctx, client
}
