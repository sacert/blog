#!/bin/bash

# Check if Air is installed
if ! command -v air &> /dev/null; then
    echo "Air is not installed. Installing now..."
    go install github.com/cosmtrek/air@latest
fi

# Run Air for live reloading
air
