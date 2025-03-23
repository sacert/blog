#!/bin/bash

# This script prepares a fresh DigitalOcean Droplet for your Go Markdown Blog deployment
# Run this script once on your new Droplet before the first GitHub Actions deployment

# Update package lists
echo "Updating package lists..."
apt-get update

# Install required packages
echo "Installing Docker and other dependencies..."
apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -

# Add Docker repository
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

# Update package lists again
apt-get update

# Install Docker
apt-get install -y docker-ce docker-ce-cli containerd.io

# Enable Docker service
systemctl enable docker
systemctl start docker

# Create directory for blog content
echo "Creating directory for blog content..."
mkdir -p /opt/blog/content

# Set permissions
echo "Setting permissions..."
chmod -R 755 /opt/blog

echo "Setup complete! Your Droplet is now ready for deployment."
echo "Don't forget to add the GitHub secrets required for the workflow:"
echo "  - SSH_PRIVATE_KEY: Your private SSH key for connecting to this Droplet"
echo "  - SSH_KNOWN_HOSTS: The SSH known hosts entry for this Droplet"
echo "  - DO_USER: The username used to connect (usually 'root')"
echo "  - DO_HOST: The IP address or hostname of this Droplet"
