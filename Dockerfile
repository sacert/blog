FROM ubuntu:22.04

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apt-get update && apt-get install -y \
    golang \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o blog-app .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./blog-app"]
