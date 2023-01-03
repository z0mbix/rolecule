set shell := ["bash", "-uc"]

# Show available targets/recipes
default:
    @just --choose

# Clean up old files
clean:
    rm -rf ./dist/*
    rm ./rolecule

# Build the binary for the current os/arch
build:
    go build -o bin/rolecule

# Configure your host to use this repo
setup:
    direnv allow

# Show git tags
tags:
    @git tag | sort -V

# Run unit tests
test:
    go test -v ./...

help:
    @just --list --list-prefix '  ❯ '
