package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create a server with timeouts
	server := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server started at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
