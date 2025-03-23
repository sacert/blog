## Deployment to DigitalOcean

This project includes a GitHub Actions workflow to deploy to a DigitalOcean Droplet.

### Initial Droplet Setup

1. Create a new DigitalOcean Droplet with Ubuntu
2. SSH into your Droplet
3. Upload and run the `droplet-setup.sh` script:

```bash
# On your local machine
scp droplet-setup.sh user@your-droplet-ip:/tmp/

# On your Droplet
chmod +x /tmp/droplet-setup.sh
sudo /tmp/droplet-setup.sh
```

### Setting Up GitHub Actions

Add the following secrets to your GitHub repository:

- `SSH_PRIVATE_KEY`: Your private SSH key for connecting to the Droplet
- `SSH_KNOWN_HOSTS`: The SSH known hosts entry for your Droplet (run `ssh-keyscan your-droplet-ip` to get this)
- `DO_USER`: The username used to connect (usually 'root')
- `DO_HOST`: The IP address or hostname of your Droplet

### Manual Deployment

If needed, you can also manually trigger the deployment from the GitHub Actions tab in your repository.

# Go Markdown Blog

A simple blog built with Go that uses Markdown files for blog posts. Features tag support and Docker deployment.

## Features

- 📝 Markdown content for easy writing
- 🏠 Home page listing all posts
- 📄 Individual post pages
- 🏷️ Tag support for categorizing posts
- 🐳 Docker support for easy deployment
- 🛠️ GitHub Actions workflow for automated deployment
- 🎨 Responsive design
- 🚀 Fast and lightweight

### Running with Docker

The blog can also be run using Docker:

```bash
# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

Alternatively, you can build and run the Docker container manually:

```bash
# Build the Docker image
docker build -t go-markdown-blog .

# Run the container
docker run -p 8080:8080 -v $(pwd)/content:/app/content go-markdown-blog
```

The blog will be available at http://localhost:8080

## Project Structure

```
blog/
├── .github/           # GitHub related files
│   └── workflows/      # GitHub Actions workflows
│       └── deploy.yml   # Deployment workflow
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
├── .dockerignore      # Files to exclude from Docker build
├── .gitignore         # Git ignore file
├── dev.sh             # Script to run development server with live reloading
├── docker-compose.yml # Docker Compose configuration
├── Dockerfile         # Docker build instructions
├── droplet-setup.sh   # Script to prepare DigitalOcean Droplet
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
- [Docker](https://www.docker.com/) - Container platform (optional)
- [GitHub Actions](https://github.com/features/actions) - CI/CD platform (optional)

## License

MIT
