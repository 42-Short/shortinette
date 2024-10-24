# shortinette: Automated Grading Framework for Coding Bootcamps
🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧
## Overview

shortinette is a framework designed to manage and automate the grading process for coding bootcamps, which we call `Shorts`. This system provides tools for efficiently and safely running and evaluating student submissions.

## Key Features

- Automated grading triggered by GitHub webhooks
- Secure execution of untrusted code using Docker
- Easily extensible

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

# All dependencies required to build the Rust modules
RUN apt-get update && apt-get install -y curl build-essential sudo m4

# Install Go
RUN curl -OL https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm go1.22.5.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Add 'student' user for running tests without permissions (default user in containers is root)
RUN useradd -m student
RUN chmod 777 /root
USER student 
RUN rustup default stable
USER root

RUN echo 'export PATH=$PATH:/root/.cargo/bin' >> /etc/profile.d/rust_path.sh

# Install 'cargo-valgrind' for testing leaks
RUN apt-get install -y valgrind
RUN /root/.cargo/bin/cargo install cargo-valgrind

WORKDIR /app

COPY ./internal /app/internal
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go

RUN go build .
```
The Docker image must be named `shortinette-testenv` (`docker build -t shortinette-testenv .` in the Dockerfile's directory). Configurability of the image's name will be added in a future release.

Note: The only mounted directory will be `./traces`, since all containers need to write into the same trace file. Everything else is just copied into the container by `shortinette` to prevent code execution on the host.

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
        "ex00",
        []string{"main.rs"},
        10,
        helloWorldTest,
    )
}
```

shortinette stores the directory the repository got cloned into into the `CloneDirectory` field of the `Exercise` struct it passes to your testing function. This way, you can add more extensive tests like lints, etc by accessing the repo's content directly.

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
WEBHOOK_URL="<Host>:8080/webhook"

# This is the organization under which the repositories will be created.
GITHUB_ORGANISATION="Your GitHub organization's name"

CONFIG_PATH="Path to your Short config.json, see below for details"

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

