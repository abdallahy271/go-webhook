package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// GitHub repository details
const (
	// sourceOwner = "abdallahy271"
	// sourceRepo       = "webhooks-test"
	targetOwner = "abdallahy271"
	targetRepo  = "go-webhook"
	// sourceBranchName = "webhook3"
	targetBranchName = "main"
)

// File details
const (
// username    = ""
// fileContent = "New content for the file"

// fileOwner = "CS404-Startup"
// fileRepo  = "Pigeon"
// filePath  = "docker-compose.yml"
)

func CreatePR(changeInfo *ChangeInfo) error {
	// Create HTTP client with authorization header
	client := &http.Client{}

	// Create pull request payload
	payload := map[string]interface{}{
		"title": "Update file",
		"body":  "Updating file content",
		"head":  fmt.Sprintf("%s:%s", targetOwner, changeInfo.SourceBranch),
		"base":  targetBranchName,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal payload:", err)
		return err
	}

	// Make a POST request for a pull request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", targetOwner, targetRepo), bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return err
	}
	accessToken, _ := os.LookupEnv("GH_ACCESS_TOKEN")

	req.Header.Set("Authorization", "token "+accessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return err
	}

	// Check response status code
	switch resp.StatusCode {
	case http.StatusCreated:
		fmt.Println("Pull request created successfully!")
	case http.StatusUnprocessableEntity:
		fmt.Println("A pull request with the same head already exists. Skipping.")
	default:
		fmt.Println("Request failed with status code:", resp.StatusCode)
		fmt.Println("Response:", string(body))
		return err
	}
	return nil
}
