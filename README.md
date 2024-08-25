# shortinette: Automated Grading Framework for Coding Bootcamps

## Overview

Shortinette is a comprehensive framework designed to manage and automate the grading process for coding bootcamps, specifically tailored for "Shorts" (coding courses). This system provides a suite of tools for efficiently running and evaluating student submissions.

## Key Features

- Automated grading triggered by GitHub webhooks
- Secure execution of untrusted code using Docker
- Comprehensive logging and error reporting
- Modular architecture for easy extension and maintenance

## Architecture

shortinette is composed of several interconnected packages, each responsible for a specific aspect of the grading pipeline:

1. **logger**: Manages logging throughout the framework, capturing important events and errors.
2. **requirements**: Validates environment variables and dependencies.
3. **testutils**: Provides utilities for compiling and running code submissions.
4. **db**: Handles interactions with the SQLite database.
5. **git**: Manages GitHub interactions, including repository management and file operations.
6. **exercise**: Defines the structure and behavior of individual coding exercises.
7. **module**: Organizes exercises into cohesive curriculum modules.
8. **webhook**: Enables automated grading triggered by GitHub events.
9. **short**: Orchestrates the entire grading process, integrating all sub-packages.

## Implementation Guide

### Prerequisites

- Docker
- Go programming environment
- GitHub account and personal access token
- Public IP for the GitHub webhook

### Step 1: Prepare the Docker Environment

Create a Dockerfile that includes all necessary dependencies for both Shortinette and the programming language being tested. Example for Rust:

```dockerfile
FROM debian:latest

# Install essential tools and Rust
RUN apt-get update && apt-get install -y curl build-essential && \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Install Go
RUN curl -OL https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm go1.22.5.linux-amd64.tar.gz

# Set up environment
ENV PATH="/usr/local/go/bin:/root/.cargo/bin:${PATH}"
WORKDIR /app
COPY . .
RUN go build .
```
The Docker image must be named `shortinette-testenv` (`docker build -t shortinette-testenv .` in the Dockerfile's directory). Configurability of the image's name will be added in a future release.

### Step 2: Define an Exercise

Create a Go file to define your exercise:

```go
package main

import (
    "github.com/42-Short/shortinette/pkg/interfaces/exercise"
    "github.com/42-Short/shortinette/pkg/testutils"
)

func helloWorldTest(ex *exercise.Exercise) exercise.Result {
    // Implement test logic here
}

func createExampleExercise() exercise.Exercise {
    return exercise.NewExercise(
        "example",
        "studentcode",
        "ex00",
        []string{"main.rs"},
        nil,
        10,
        helloWorldTest,
    )
}
```

### Step 3: Define a Module

Create a module that includes your exercise:

```go
package main

import (
    "github.com/42-Short/shortinette/pkg/interfaces/module"
    "github.com/42-Short/shortinette/pkg/logger"
)

func createExampleModule() module.Module {
    exercises := map[string]exercise.Exercise{
        "hello-world": createExampleExercise(),
    }

    module, err := module.NewModule(
        "example-module",
        50,
        exercises,
        "subjects/ex00",
    )
    if err != nil {
        logger.Error.Fatalf("Failed to create module: %v", err)
    }
    return module
}
```

### Step 4: Initialize and Run Shortinette

Set up the main function to run Shortinette:

```go
package main

import (
    "github.com/42-Short/shortinette/pkg/short"
    "github.com/42-Short/shortinette/pkg/webhook"
)

func main() {
    modules := map[string]module.Module{
        "example-module": createExampleModule(),
    }

    testMode := webhook.NewWebhookTestMode(modules)
    s := short.NewShort("Example Shortinette", modules, testMode)
    s.Start("example-module")
}
```

### Step 5: Configure the Environment
#### .env File
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
#### Configuration File
Now configure the .json file whose path you set in your environment:
```json
{
    "participants": [
        {
            "github_username": "shortinette-test",
            "intra_login": "1"
        }
    ]
}
```
_note: The intra_login variable is used to build the names of the repos which will be created. You can of course set it to something else if you want the repos to be named differently. The repo naming format is: <intra_login>-<module_name>_

### Step 6: Run Shortinette

Execute the Shortinette system using the command:

```
go run .
```

For more advanced implementations and examples, refer to the [rust piscine repository](https://github.com/42-Short/rust).

