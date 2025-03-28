package handlers

import (
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/sacert/blog/models"
)

// BlogHandler contains all the dependencies needed for our handlers
type BlogHandler struct {
	ContentDir string
	Templates  map[string]*template.Template
}

// NewBlogHandler creates a new BlogHandler instance
func NewBlogHandler() *BlogHandler {
	tmpl := make(map[string]*template.Template)
	tmpl["home"] = template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
	tmpl["post"] = template.Must(template.ParseFiles("templates/base.html", "templates/post.html"))

	return &BlogHandler{
		ContentDir: "content",
		Templates:  tmpl,
	}
}

// ListPosts handles the homepage showing all posts
func (h *BlogHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	posts, err := models.GetPosts(h.ContentDir)
	if err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishDate.After(posts[j].PublishDate)
	})

	data := models.BlogData{
		Posts:       posts,
		Title:       "My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
	}

	err = h.Templates["home"].ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ShowPost handles displaying a single post
func (h *BlogHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/post/")

	posts, err := models.GetPosts(h.ContentDir)
	if err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}

	var post models.Post
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

	data := models.BlogData{
		Posts:       []models.Post{post},
		Title:       post.Title + " - My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
	}

	err = h.Templates["post"].ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
