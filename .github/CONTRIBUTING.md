# Contributing to shortinette
We want to make contributing to this project as easy and transparent as possible, whether it's:
*  Reporting a bug
*  Discussing the current state of the code
*  Submitting a fix
*  Proposing new features

## We Develop with GitHub
We use GitHub to host code, track issues and feature requests, as well as accept pull requests.
We use [GitHub FLow](https://guides.github.com/introduction/flow/index.html): **All Changes Happen Through Pull Requests**.

### Report Bugs using GitHub's Issues
We use GitHub issues to track bugs. Report a bug by [opening an issue](https://github.com/42-Short/shortinette/issues); it's very easy!

### Write Bug Reports with Detail, Background and Sample Code!
**Good Bug Reports** tend to have:
*  A quick summary and/or background
*  Steps to reproduce
   * Be specific!
   * Give sample code if you can.
*  What you expected would happen.
*  What actually happens.
*  Notes, for example:
   * What you think might be happening.
   * Stuff you tried that did not work.

The more thorough the bug report, the easier the fix :)

### Development Workflow
1. Create a branch `git checkout -b your-github-username/(feat|fix|docs)/branch-name`
2. Make your changes.
3. **Test your changes!**
4. Commit your changes: `git commit -m "(add|test|chore|style|fix): description-of-your-commit`
5. Push to your branch.
6. Open a Pull Request! Add screenshots, explanation of your thought process, anything that could be useful for code review.

### Code Review
Pull requests are reviewed by the project maintainers. Here's what we look for:
* Code quality
* Adequate testing coverage
* Clear and concise commit messages
* Relevance and necessity of the proposed changes.
* Documentation: If anything is unclear about what has been changed, how to test, (...), we _will_ deny the pull request.

### Updating Documentation

If your changes affect the project documentation (e.g. new features, changes to existing functionality), please update the relevant documentation files as part of your pull request. This helps keep our documentation up to date and useful for everyone.

### Coding Style
We use the Golang [default formatter](https://github.com/golang/vscode-go) with default configuration. Please also be careful to write small, single-purpose functions. This will keep the code easy to read and maintain.
