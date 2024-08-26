// Package shortinette is the core framework for managing and automating the process of
// grading coding bootcamps (Shorts). It provides a comprehensive set of tools for
// running and testing student submissions across various programming languages.
// The shortinette package is composed of several sub-packages, each responsible for a specific
// aspect of the grading pipeline:
//
//   - `logger`: Handles logging for the framework, including general informational messages,
//     error reporting, and trace logging for feedback on individual submissions. This package ensures
//     that all important events and errors are captured for debugging and auditing purposes.
//
//   - `requirements`: Validates the necessary environment variables and dependencies required
//     by the framework. This includes checking for essential configuration values in a `.env` file
//     and ensuring that all necessary tools (e.g., Docker images) are available before grading begins.
//
//   - `testutils`: Provides utility functions for compiling and running code submissions.
//     This includes functions for compiling Rust code, running executables with various
//     options (such as timeouts and real-time output), and manipulating files. The utility
//     functions are designed to handle the intricacies of running untrusted student code
//     safely and efficiently.
//
//   - `git`: Manages interactions with GitHub, including cloning repositories, managing
//     collaborators, and uploading files. This package abstracts the GitHub API to simplify
//     common tasks such as adding collaborators to repositories, creating branches, and
//     pushing code or data to specific locations in a repository.
//
//   - `exercise`: Defines the structure and behavior of individual coding exercises.
//     This includes specifying the files that students are allowed to submit, the expected
//     output, and the functions to be tested. The `exercise` package provides the framework
//     for setting up exercises, running tests, and reporting results.
//
//   - `module`: Organizes exercises into modules, allowing for the grouping of related exercises
//     into a coherent curriculum. The `module` package handles the execution of all exercises
//     within a module, aggregates results, and manages the overall grading process.
//
//   - `webhook`: Enables automatic grading triggered by GitHub webhooks. This allows for a
//     fully automated workflow where student submissions are graded as soon as they are
//     pushed to a specific branch in a GitHub repository.
//
//   - `short`: The central orchestrator of the grading process, integrating all sub-packages
//     into a cohesive system. The `short` package handles the setup and teardown of grading
//     environments, manages the execution of modules and exercises, and ensures that all
//     results are properly recorded and reported.
// package shortinette

package main

import (
	"github.com/42-Short/shortinette/internal/tests/R00"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/short"
	"github.com/42-Short/shortinette/pkg/short/testmodes/webhook"
)

func main() {
	modules := map[string]Module.Module{
		"00": *R00.R00(),
		// TODO: "01": *R01.R01(), // TODO
		// TODO: "02": *R02.R02(), // TODO
		// TODO: "03": *R03.R03(), // TODO
		// TODO: "04": *R04.R04(), // TODO
		// TODO: "05": *R05.R05(), // TODO
	}
	short := short.NewShort("Rust Piscine 1.0", modules, webhook.NewWebhookTestMode(modules))
	short.Start("00")
}
