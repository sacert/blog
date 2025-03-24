package main

import (
	_ "fmt" // Added for fmt.Println in the markdown example
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sacert/blog/handlers"
	"github.com/sacert/blog/models"
)

func BenchmarkListPosts(b *testing.B) {
	blogHandler := handlers.NewBlogHandler()

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatal(err)
	}

	// Reset the timer to exclude setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create a ResponseRecorder for each iteration
		rr := httptest.NewRecorder()

		// Call the handler
		blogHandler.ListPosts(rr, req)

		// Check for successful response
		if status := rr.Code; status != http.StatusOK {
			b.Fatalf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}
}

func BenchmarkMdToHTML(b *testing.B) {
	markdown := `
# Test Heading

This is a paragraph with **bold** and *italic* text.

## Subheading

- List item 1
- List item 2
- List item 3

~~~go
func example() {
	fmt.Println("Hello, world!")
}
~~~

> This is a blockquote.

[Link text](https://example.com)

![Image alt text](image.jpg)
`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		models.MdToHTML(markdown)
	}
}

func BenchmarkParsePost(b *testing.B) {
	// Save the original function and restore it after the test
	origGetPosts := models.GetPosts
	defer func() { models.GetPosts = origGetPosts }()

	// Since we can't easily mock file operations in benchmarks,
	// we'll create a minimal test scenario
	const contentDir = "models/testdata"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Don't need to cast the function back to itself - remove the unnecessary conversion
		models.GetPosts = origGetPosts

		posts, err := models.GetPosts(contentDir)
		if err != nil {
			b.Fatalf("GetPosts() error = %v", err)
		}

		if len(posts) == 0 {
			b.Fatalf("Expected posts to be returned")
		}
	}
}
