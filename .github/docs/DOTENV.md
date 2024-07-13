# Documentation for `.env` File Configuration
Create a `.env` file at the root of your repository and fill it up like below:
```.env
# These are used for identifying you when making requests on the GitHub API.
GITHUB_ADMIN="Your GitHub username"
GITHUB_EMAIL="Your GitHub email"
GITHUB_TOKEN="Your GitHub personal access token"

# We use Webhooks to record events on repositories.
WEBHOOK_URL="<HOST>:8080/webhook"

# This is the organization under which the repositories will be created.
GITHUB_ORGANISATION="Your GitHub organization's name"

CONFIG_PATH="Path to your Short config"

```
