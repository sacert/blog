# Getting Started with Golang

Tags: golang, tutorial, programming

Go, also known as Golang, is an open-source programming language designed at Google. Let's explore how to get started with Go development.

## Installation

First, download and install Go from the [official website](https://golang.org/dl/). After installation, verify it's working by running:

```bash
go version
```

## Setting Up Your Workspace

Go projects are typically structured with a specific workspace layout:

```
myproject/
├── go.mod
├── main.go
└── pkg/
    └── mypackage/
        └── mypackage.go
```

Initialize a new module:

```bash
mkdir myproject
cd myproject
go mod init example.com/myproject
```

## Your First Go Program

Create a file named `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

Run it with:

```bash
go run main.go
```

## Building and Installing

To compile your program:

```bash
go build
```

This creates an executable in your current directory. To install it to your Go bin directory:

```bash
go install
```

Now you're ready to start building with Go!