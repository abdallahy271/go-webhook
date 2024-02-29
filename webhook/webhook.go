package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GitHub repository details
const (
	sourceOwner      = "abdallahy271"
	sourceRepo       = "webhooks-test"
	targetOwner      = "abdallahy271"
	targetRepo       = "go-webhook"
	sourceBranchName = "webhook"
	targetBranchName = "main"
)

// Personal access token for authentication
const accessToken = "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"

// File details
// const (
// 	filePath    = "path/to/file.txt"
// 	fileContent = "New content for the file"
// )

func CreatePR() {
	if err := CommitChange(); err != nil {
		return
	}

	// Create pull request payload
	payload := map[string]interface{}{
		"title": "Update file",
		"body":  "Updating file content",
		"head":  fmt.Sprintf("%s:%s", sourceOwner, sourceBranchName),
		"base":  targetBranchName,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal payload:", err)
		return
	}

	// Create HTTP client with authorization header
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", targetOwner, targetRepo), bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("Authorization", "token "+accessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	// Check response status code
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Request failed with status code:", resp.StatusCode)
		fmt.Println("Response:", string(body))
		return
	}

	fmt.Println("Pull request created successfully")
}
