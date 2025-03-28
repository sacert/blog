package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sacert/blog/models"
)

func setupTestTemplates(t *testing.T) map[string]*template.Template {
	// Create temporary test templates
	tmpDir, err := os.MkdirTemp("", "blogtest")
	if err != nil {
		t.Fatalf("Could not create temp dir: %v", err)
	}

	// Clean up after the test
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to clean up temp dir: %v", err)
		}
	})

	// Create the template files
	baseHTML := `{{ define "base" }}<!DOCTYPE html><html><head><title>{{.Title}}</title></head><body>{{template "content" .}}</body></html>{{ end }}`
	homeHTML := `{{ define "content" }}<h1>Posts</h1><div class="posts">{{ range .Posts }}<div class="post"><h2>{{ .Title }}</h2></div>{{ end }}</div>{{ end }}`
	postHTML := `{{ define "content" }}{{ range .Posts }}<article><h1>{{ .Title }}</h1><div>{{ .Content }}</div></article>{{ end }}{{ end }}`

	// Write the template files
	if err := os.WriteFile(filepath.Join(tmpDir, "base.html"), []byte(baseHTML), 0o644); err != nil {
		t.Fatalf("Could not write base template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "home.html"), []byte(homeHTML), 0o644); err != nil {
		t.Fatalf("Could not write home template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "post.html"), []byte(postHTML), 0o644); err != nil {
		t.Fatalf("Could not write post template: %v", err)
	}

	// Parse the templates
	tmpl := make(map[string]*template.Template)
	tmpl["home"] = template.Must(template.ParseFiles(filepath.Join(tmpDir, "base.html"), filepath.Join(tmpDir, "home.html")))
	tmpl["post"] = template.Must(template.ParseFiles(filepath.Join(tmpDir, "base.html"), filepath.Join(tmpDir, "post.html")))

	return tmpl
}

func createMockPosts() []models.Post {
	return []models.Post{
		{
			Title:       "Test Post 1",
			Slug:        "test-post-1",
			Content:     template.HTML("<p>Test content 1</p>"),
			RawContent:  "Test content 1",
			PublishDate: time.Now(),
			Summary:     "Test summary 1",
		},
		{
			Title:       "Test Post 2",
			Slug:        "test-post-2",
			Content:     template.HTML("<p>Test content 2</p>"),
			RawContent:  "Test content 2",
			PublishDate: time.Now().Add(-24 * time.Hour),
			Summary:     "Test summary 2",
		},
	}
}

func TestListPosts(t *testing.T) {
	// Save the original function and restore it after the test
	origGetPosts := models.GetPosts
	defer func() { models.GetPosts = origGetPosts }()

	// Override GetPosts to use mock data
	models.GetPosts = func(string) ([]models.Post, error) {
		return createMockPosts(), nil
	}

	tmpl := setupTestTemplates(t)

	h := &BlogHandler{
		ContentDir: "testdir",
		Templates:  tmpl,
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ListPosts)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the response body contains the post titles
	if body := rr.Body.String(); !contains(body, "Test Post 1") || !contains(body, "Test Post 2") {
		t.Errorf("handler returned unexpected body: got %v", body)
	}
}

func TestShowPost(t *testing.T) {
	// Save the original function and restore it after the test
	origGetPosts := models.GetPosts
	defer func() { models.GetPosts = origGetPosts }()

	// Override GetPosts to use mock data
	models.GetPosts = func(string) ([]models.Post, error) {
		return createMockPosts(), nil
	}

	tmpl := setupTestTemplates(t)

	h := &BlogHandler{
		ContentDir: "testdir",
		Templates:  tmpl,
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/post/test-post-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ShowPost)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the response body contains the post title
	if body := rr.Body.String(); !contains(body, "Test Post 1") {
		t.Errorf("handler returned unexpected body: got %v", body)
	}

	// Test non-existent post
	req, err = http.NewRequest("GET", "/post/non-existent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code - should be Not Found
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
