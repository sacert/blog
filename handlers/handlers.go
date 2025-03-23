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

	// Get all unique tags
	allTags := models.GetAllTags(posts)

	data := models.BlogData{
		Posts:       posts,
		Title:       "My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
		AllTags:     allTags,
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

	allTags := models.GetAllTags(posts)

	data := models.BlogData{
		Posts:       []models.Post{post},
		Title:       post.Title + " - My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
		AllTags:     allTags,
	}

	err = h.Templates["post"].ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleTag handles filtering posts by tag
func (h *BlogHandler) HandleTag(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimPrefix(r.URL.Path, "/tag/")
	tag = strings.TrimSpace(tag)

	if tag == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	posts, err := models.GetPosts(h.ContentDir)
	if err != nil {
		http.Error(w, "Error reading posts", http.StatusInternalServerError)
		return
	}

	// Filter posts by tag
	var filteredPosts []models.Post
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

	allTags := models.GetAllTags(posts)

	data := models.BlogData{
		Posts:       filteredPosts,
		Title:       "Posts tagged '" + tag + "' - My Go Markdown Blog",
		CurrentYear: time.Now().Year(),
		AllTags:     allTags,
		ActiveTag:   tag,
	}

	err = h.Templates["home"].ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
