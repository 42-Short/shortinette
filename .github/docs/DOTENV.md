# Documentation for `.env` File Configuration
The `shortinette` backend needs to interact with the GitHub API in order to create/pull from repositories.
Therefore, you need to provide your credentials.

1. Create a `.env` file at the root of your repository.
2. This is how you need to fill it up:
```.env
# These are used for identifying you when making requests on the GitHub API.
GITHUB_USER="your GitHub username"
GITHUB_TOKEN="your GitHub personal access token"

# This is the organization under which the repositories will be created. Please create your own for testing purposes.
GITHUB_ORGANISATION="your GitHub organization's name"
```
