# Blog

## Project Structure

```
blog/
├── content/           # Markdown blog posts
├── static/            # Static assets (CSS, JS, images)
│   └── css/
│       └── style.css
├── templates/         # HTML templates
│   ├── base.html
│   ├── home.html
│   └── post.html
├── tmp/               # Temporary directory for Air (live reloading)
├── .air.toml          # Configuration for Air (live reloading)
├── .gitignore         # Git ignore file
├── dev.sh             # Script to run development server with live reloading
├── go.mod             # Go module file
├── main.go            # Main application file
└── README.md          # This file
```

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone this repository
2. Install dependencies:

```bash
go mod download
```

### Running the Blog

#### Standard Mode

```bash
go run main.go
```

#### Development Mode with Live Reloading

```bash
# Make the script executable if needed
chmod +x dev.sh

# Run the development server with live reloading
./dev.sh
```

The blog will be available at http://localhost:8080 in both modes. In development mode, the server will automatically reload when you make changes to your files.

### Creating Blog Posts

Create new `.md` files in the `content` directory. The first line of each file should be a level 1 heading (`# Title`) which will be used as the post title.

Example:

```markdown
# My New Post

This is the content of my new post.

## Subheading

More content here...
```

The filename (without the `.md` extension) will be used as the URL slug for the post.

## Customization

- Edit the templates in the `templates` directory to change the HTML structure
- Modify `static/css/style.css` to change the appearance
- Update the blog title and other settings in `main.go`

## Dependencies

- [gomarkdown/markdown](https://github.com/gomarkdown/markdown) - Markdown to HTML conversion
- [cosmtrek/air](https://github.com/cosmtrek/air) - Live reload for Go apps (development only)

## License

MIT
