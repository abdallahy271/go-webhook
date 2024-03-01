package main

import (
	"encoding/json"
	"fmt"
	"go-webhook/webhook"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/joho/godotenv"
)

// Define the structure of the webhook payload
type WebhookPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Merged bool   `json:"merged"`
		URL    string `json:"url"`
	} `json:"pull_request"`
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Define a simple handler function
	handler := func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is an OPTIONS preflight request
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method == "POST" {

			fmt.Println("received request")

			// Decode the JSON payload
			var payload WebhookPayload
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				http.Error(w, "Failed to parse webhook payload", http.StatusBadRequest)
				return
			}

			// Check if the action is "closed" and the pull request was merged
			if payload.Action == "closed" && payload.PullRequest.Merged {
				fmt.Println("Received merged pull request webhook", payload.PullRequest.URL)
				// Here you can perform actions specific to merged pull requests
				// Fetch the changes that were merged in the pull request
				changeInfo, err := webhook.GetMergedPullRequestChanges(payload.PullRequest.URL)
				if err != nil {
					fmt.Printf("Error fetching changes for merged pull request: %v\n", err)
				} else {
					fmt.Println("Changes that were merged:")
					for _, change := range changeInfo.Changes {
						fmt.Println(change)
						webhook.CreatePR(change)
					}

				}
			} else {
				// Ignore the webhook if it's not a closed pull request or if it wasn't merged
				fmt.Println("Ignoring webhook")
			}
		}

		// webhook.CreatePR()
		// ExecuteCommand()

		io.WriteString(w, "Hello, World!\n") // Respond with "Hello, World!"
	}

	// Register the handler function to handle all requests to the root URL "/"
	http.HandleFunc("/", handler)

	// Start the server on port 3333
	fmt.Println("Server started on http://localhost:3333")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func ExecuteCommand() {
	// Command to execute
	cmd := "openapi-changes"
	// Arguments for the command (if any)
	args := []string{}

	// Create a new Cmd struct to represent the command
	command := exec.Command(cmd, args...)

	// Run the command and capture its output
	output, err := command.Output()
	if err != nil {
		log.Fatalf("Error running command: %s", err)
	}

	// Print the output
	fmt.Println(string(output))
}
