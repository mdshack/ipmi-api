# GitHub Workflows

This directory contains the GitHub Actions workflows for this project:

## Workflows

1. **`release-please.yml`** - Automatically creates release PRs and releases when changes are merged to main. Also builds and publishes Docker images to both Docker Hub and GitHub Container Registry on release.
2. **`test.yml`** - Runs Go tests on pull requests and pushes to main

## Secrets Required

For the Docker Hub workflow to work, you need to set the following secrets in your repository:

- `DOCKER_USERNAME` - Your Docker Hub username
- `DOCKER_PASSWORD` - Your Docker Hub password or access token

The release-please workflow uses `GITHUB_TOKEN` which is automatically provided by GitHub Actions.