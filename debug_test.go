package main

import (
	"fmt"
	"testing"

	"github.com/sacert/blog/models"
)

func TestDebugMdToHTML(t *testing.T) {
	// This is the exact same input as in the code block test
	markdown := "```go\nfunc test() {\n  fmt.Println(\"hello\")\n}\n```"
	output := models.MdToHTML(markdown)
	
	// Print the exact output for comparison
	fmt.Println("Actual output (for use in test):")
	fmt.Printf("%q\n", output)
}
