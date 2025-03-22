package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Title       string
	Slug        string
	Content     template.HTML
	RawContent  string
	PublishDate time.Time
	Summary     string
}

type BlogData struct {
	Posts []Post
	Title string
	CurrentYear int
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Home page route
	http.HandleFunc("/", listPostsHandler)

	// Individual post route
	http.HandleFunc("/post/", postHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := getPosts()
	if err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishDate.After(posts[j].PublishDate)
	})

	data := BlogData{
		Posts: posts,
		Title: "My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")
	
	posts, err := getPosts()
	if err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}

	var post Post
	found := false
	for _, p := range posts {
		if p.Slug == slug {
			post = p
			found = true
			break
		}
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	data := BlogData{
		Posts: []Post{post},
		Title: post.Title + " - My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
	}

	tmpl, err := template.ParseFiles("templates/base.html", "templates/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPosts() ([]Post, error) {
	var posts []Post
	files, err := os.ReadDir("content")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		slug := strings.TrimSuffix(file.Name(), ".md")
		filePath := filepath.Join("content", file.Name())
		
		// Get file info to access ModTime
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}
		
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// Extract title from first line
		lines := strings.SplitN(string(content), "\n", 3)
		title := strings.TrimPrefix(lines[0], "# ")
		
		// Parse content to exclude title
		contentWithoutTitle := strings.Join(lines[1:], "\n")
		
		// Create summary (first 150 chars)
		summary := contentWithoutTitle
		if len(summary) > 150 {
			summary = summary[:150] + "..."
		}

		// Convert markdown to HTML
		htmlContent := mdToHTML(contentWithoutTitle)

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

func mdToHTML(md string) string {
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
