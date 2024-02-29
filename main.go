package main

import (
	"fmt"
	"io"
	"net/http"
)

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

		fmt.Println("received request")
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
