package webhook

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ChangeInfo struct {
	Changes      []string
	Owner        string
	Repo         string
	SourceBranch string
}

func GetMergedPullRequestChanges(prURL string) (*ChangeInfo, error) {
	// Initialize GitHub client with authentication
	ctx, client := GetGitHubClient()

	// Extract owner and repository name from PR URL
	owner, repo, err := extractOwnerAndRepoFromURL(prURL)
	if err != nil {
		return nil, err
	}

	// Extract PR number from PR URL
	prNumber, err := extractPRNumberFromURL(prURL)
	if err != nil {
		return nil, err
	}

	// Get the files changed in the pull request
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return nil, err
	}

	var changes []string
	for _, file := range files {
		changes = append(changes, *file.Filename)
	}
	currentTimeMilli := time.Now().UnixNano() / int64(time.Millisecond)
	sourceBranch := fmt.Sprintf("%s-%d", owner, currentTimeMilli)

	changeInfo := &ChangeInfo{
		Changes:      changes,
		Owner:        owner,
		Repo:         repo,
		SourceBranch: sourceBranch,
	}

	return changeInfo, nil
}

func extractOwnerAndRepoFromURL(url string) (string, string, error) {
	// Extract owner and repository name from URL (e.g., https://github.com/owner/repo/pull/123)
	// You may need to adjust this function based on the format of your URL
	// For GitHub, the URL format is typically: https://github.com/owner/repo/pull/123
	// You may need to handle other formats depending on your version control system
	// Here, we assume that the URL format is consistent with GitHub's format
	// You may need to implement different logic for other version control systems
	// For example, for GitLab, the URL format might be different
	// You would need to adapt this function accordingly
	// This is a simplified example assuming the URL format is consistent
	// You might need to handle edge cases and error conditions
	// For a production-ready implementation, consider using a more robust URL parsing library
	// and handle errors and edge cases appropriately
	// Here, we just split the URL by "/" and extract the owner and repository name
	parts := strings.Split(url, "/")
	if len(parts) < 7 {
		return "", "", fmt.Errorf("invalid URL format: %s", url)
	}
	return parts[4], parts[5], nil
}

func extractPRNumberFromURL(url string) (int, error) {
	// Extract PR number from URL (e.g., https://github.com/owner/repo/pull/123)
	// Similar to extractOwnerAndRepoFromURL function, you may need to adjust this function
	parts := strings.Split(url, "/")
	if len(parts) < 7 {
		return 0, fmt.Errorf("invalid URL format: %s", url)
	}
	prNumberStr := parts[7]
	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PR number: %s", prNumberStr)
	}
	return prNumber, nil
}
