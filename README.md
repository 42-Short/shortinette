# Shortinette

Welcome to the official repository for Shortinette! This project is part of a plan to organize a `Rust` Piscine at 42Vienna. 

## Table of Contents
- [About the Project](#about-the-project)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About the Project

Shortinette is the backend of the application which will be responsible for creating teams and checking student's submissions. It currently includes creation of repos for submissions, check for forbidden items (functions, macros, keywords) and automated unit-testing.

## Tech Stack

The application will be containerized and built with Docker Compose (https://github.com/42-Short/shortinette/issues/7).

- **Backend**: Go
- **Frontend**: TBD (Simple Web-App or CLI)

## Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

Make sure you have the following installed:
- Go

### Usage

1. Clone the repo:
```sh
git clone https://github.com/42-Short/shortinette.git
cd website
```
2. Create your .env file and add the necessary variables (See [DOTENV](.github/docs/DOTENV.md) for details).
3. `go run .`

## Contributing

Contributions are what make open source such an amazing place to learn, inspire and create. Any contributions are **greatly** appreciated.

See [CONTRIBUTING](.github/CONTRIBUTING.md) for more information on how to get started.

## License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.

## Contact

Click [here](https://discord.gg/WPxyu4mW) to join our development Discord server.
