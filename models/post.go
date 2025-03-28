package models

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Post represents a blog post
type Post struct {
	Title       string
	Slug        string
	Content     template.HTML // #nosec G203 - This is safe as we control the content source
	RawContent  string
	PublishDate time.Time
	Summary     string
}

// BlogData represents data passed to templates
type BlogData struct {
	Posts       []Post
	Title       string
	CurrentYear int
}

// GetPostsFunc defines the function signature for getting posts
type GetPostsFunc func(string) ([]Post, error)

// GetPosts is a variable that holds the function to get posts
// This allows for dependency injection in tests
var GetPosts GetPostsFunc = getPostsImpl

// getPostsImpl is the actual implementation of getting posts
func getPostsImpl(contentDir string) ([]Post, error) {
	var posts []Post

	// Validate the content directory first
	contentDirAbs, err := filepath.Abs(contentDir)
	if err != nil {
		return nil, err
	}

	// Check if the directory exists
	fileInfo, dirErr := os.Stat(contentDirAbs)
	if dirErr != nil || !fileInfo.IsDir() {
		return nil, dirErr
	}

	files, err := os.ReadDir(contentDirAbs)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".md")
		filePath := filepath.Join(contentDirAbs, file.Name())

		// Verify the file is still within the content directory (path traversal prevention)
		if !strings.HasPrefix(filePath, contentDirAbs) {
			continue
		}

		// Get file info to access ModTime
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		// #nosec G304 - We've validated the file path is within our content directory
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		markdownContent := string(content)

		// Extract title from first line
		lines := strings.SplitN(markdownContent, "\n", 3)
		title := strings.TrimPrefix(lines[0], "# ")

		// Parse content to exclude title
		contentWithoutTitle := strings.Join(lines[1:], "\n")
		contentWithoutTitle = strings.TrimSpace(contentWithoutTitle)

		// Create summary (first 150 chars)
		summary := contentWithoutTitle
		if len(summary) > 150 {
			summary = summary[:150] + "..."
		}

		// Convert markdown to HTML
		htmlContent := MdToHTML(contentWithoutTitle)

		// #nosec G203 - This is safe as we control the content source (markdown files)
		post := Post{
			Title:       title,
			Slug:        slug,
			Content:     template.HTML(htmlContent),
			RawContent:  contentWithoutTitle,
			PublishDate: fileInfo.ModTime(),
			Summary:     summary,
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// MdToHTML converts markdown to HTML
func MdToHTML(md string) string {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Parse markdown string
	doc := p.Parse([]byte(md))

	// Setup HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert to HTML
	return string(markdown.Render(doc, renderer))
}
