package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sacert/blog/handlers"
)

func TestIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	
	// Create a new blog handler
	blogHandler := handlers.NewBlogHandler()
	
	// Create a test HTTP server with a simple function that calls our handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fake the URL path to be / for any request
		r.URL.Path = "/"
		blogHandler.ListPosts(w, r)
	}))
	defer server.Close()
	
	// Make a request to the test server
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to GET: %v", err)
	}
	defer resp.Body.Close()
	
	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	
	// For debugging
	t.Logf("Response body: %s", string(body))
	
	// Check for expected content
	expected := "Home"
	if !strings.Contains(string(body), expected) {
		t.Errorf("Expected response to contain %q but it didn't", expected)
	}
}
