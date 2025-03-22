package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	Tags        []string
}

type BlogData struct {
	Posts      []Post
	Title      string
	CurrentYear int
	AllTags    []string
	ActiveTag  string
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Home page route
	http.HandleFunc("/", listPostsHandler)

	// Individual post route
	http.HandleFunc("/post/", postHandler)
	
	// Tag route
	http.HandleFunc("/tag/", tagHandler)

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

	// Get all unique tags
	allTags := getAllTags(posts)

	data := BlogData{
		Posts:      posts,
		Title:      "My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
		AllTags:    allTags,
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

	allTags := getAllTags(posts)

	data := BlogData{
		Posts:      []Post{post},
		Title:      post.Title + " - My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
		AllTags:    allTags,
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

func tagHandler(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimPrefix(r.URL.Path, "/tag/")
	tag = strings.TrimSpace(tag)
	
	if tag == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	posts, err := getPosts()
	if err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}
	
	// Filter posts by tag
	var filteredPosts []Post
	for _, post := range posts {
		for _, postTag := range post.Tags {
			if strings.EqualFold(postTag, tag) {
				filteredPosts = append(filteredPosts, post)
				break
			}
		}
	}
	
	// Sort posts by date (newest first)
	sort.Slice(filteredPosts, func(i, j int) bool {
		return filteredPosts[i].PublishDate.After(filteredPosts[j].PublishDate)
	})
	
	allTags := getAllTags(posts)
	
	data := BlogData{
		Posts:      filteredPosts,
		Title:      "Posts tagged '" + tag + "' - My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
		AllTags:    allTags,
		ActiveTag:  tag,
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

		markdownContent := string(content)
		
		// Extract tags if they exist (format: "Tags: tag1, tag2, tag3")
		var tags []string
		tagsRegex := regexp.MustCompile(`(?m)^Tags:\s*(.*?)$`)
		tagsMatch := tagsRegex.FindStringSubmatch(markdownContent)
		
		if len(tagsMatch) > 1 {
			// Remove the tags line from content
			markdownContent = tagsRegex.ReplaceAllString(markdownContent, "")
			
			// Process tags
			tagList := strings.Split(tagsMatch[1], ",")
			for _, tag := range tagList {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}
		
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
		htmlContent := mdToHTML(contentWithoutTitle)

		post := Post{
			Title:       title,
			Slug:        slug,
			Content:     template.HTML(htmlContent),
			RawContent:  contentWithoutTitle,
			PublishDate: fileInfo.ModTime(),
			Summary:     summary,
			Tags:        tags,
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

func getAllTags(posts []Post) []string {
	// Use a map to ensure uniqueness
	tagsMap := make(map[string]bool)
	
	for _, post := range posts {
		for _, tag := range post.Tags {
			tagsMap[tag] = true
		}
	}
	
	// Convert map keys to slice
	var tags []string
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	
	// Sort tags alphabetically
	sort.Strings(tags)
	
	return tags
}
