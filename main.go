package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sacert/blog/handlers"
)

func main() {
	// Create a new blog handler
	blogHandler := handlers.NewBlogHandler()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up routes
	http.HandleFunc("/", blogHandler.ListPosts)
	http.HandleFunc("/post/", blogHandler.ShowPost)
	http.HandleFunc("/tag/", blogHandler.HandleTag)

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
