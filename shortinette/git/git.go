// `git` is the package responsible for interactions with GitHub. The repos are _always_ assumed to be in the
// organisation specified by the `GITHUB_ORGANISATION` environment variable.
package git

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v66/github"
	"github.com/joho/godotenv"
)

func deleteRepo(name string) (err error) {
	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not delete repo '%s': %v", name, err)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	if resp, err := client.Repositories.Delete(context.Background(), orga, name); err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("could not delete repo '%s': %v", name, err)
		} else {
			fmt.Printf("repo '%s' not found in orga '%s'\n", name, orga)
			return nil
		}
	}

	fmt.Printf("repo '%s' successfully deleted\n", name)
	return nil
}

// Checks for environment variables required to interact with the GitHub API. Returns their values
// if they exist, sets the error's value if not.
func requireEnv() (githubToken string, githubOrga string, err error) {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("warning: .env file not found, this is fine in the GitHub Actions environment, this is a problem if you are running this locally\n")
	}

	missingVars := []string{}

	githubToken = os.Getenv("TOKEN_GITHUB")
	if githubToken == "" {
		missingVars = append(missingVars, "TOKEN_GITHUB")
	}

	githubOrga = os.Getenv("ORGA_GITHUB")
	if githubOrga == "" {
		missingVars = append(missingVars, "ORGA_GITHUB")
	}

	if len(missingVars) != 0 {
		err = fmt.Errorf("missing environment variable(s): %s", strings.Join(missingVars, ", "))
	}

	return githubToken, githubOrga, err
}

// Checks whether `err` is related to the repo already existing.
func isRepoAlreadyExists(err error) (exists bool) {
	if githubErr, ok := err.(*github.ErrorResponse); ok {
		for _, e := range githubErr.Errors {
			if e.Message == "name already exists on this account" {
				return true
			}
		}
	}
	return false
}

// Creates a new repository `name` under the GitHub organisation specified by the
// GITHUB_ORGANISATION environment variable. If `private` is true, the repository's
// visibility will be private.
func NewRepo(name string, private bool, description string) (err error) {
	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	repo := &github.Repository{Name: &name, Private: &private, Description: &description}

	createdRepo, response, err := client.Repositories.Create(context.Background(), orga, repo)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusUnprocessableEntity {
			if isRepoAlreadyExists(err) {
				fmt.Printf("repo %s already exists under orga %s, skipping\n", name, os.Getenv("GITHUB_ORGANISATION"))
				return nil
			}
		}
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	fmt.Printf("repo created: %s at URL: %s\n", *createdRepo.Name, *createdRepo.HTMLURL)
	return nil
}

// Adds collaborator `collaboratorName` to repo `repoName` (under the GitHub organisation specified by
// the GITHUB_ORGANISATION environment variable) with access level `permission“.
func AddCollaborator(repoName string, collaboratorName string, permission string) (err error) {
	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not add collaborator %s to repo %s: %v", collaboratorName, repoName, err)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	options := &github.RepositoryAddCollaboratorOptions{
		Permission: permission,
	}

	if _, _, err = client.Repositories.AddCollaborator(context.Background(), orga, repoName, collaboratorName, options); err != nil {
		return fmt.Errorf("could not add collaborator %s to repo %s: %v", collaboratorName, repoName, err)
	}

	fmt.Printf("user %s added to repo %s with %s access\n", collaboratorName, repoName, permission)
	return nil
}

// Clones repo `name` (from the GitHub organisation specified by the GITHUB_ORGANISATION
// environment variable). Does nothing if the directory is cloned already.
func Clone(name string) (err error) {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		fmt.Printf("'%s' seems to cloned already, returning\n", name)
		return nil
	}

	token, orga, err := requireEnv()
	if err != nil {
		return fmt.Errorf("could not clone '%s': %v", name, err)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	repo, _, err := client.Repositories.Get(context.Background(), orga, name)
	if err != nil {
		return fmt.Errorf("could not clone '%s': %v", name, err)
	}

	cmd := exec.Command("git", "clone", repo.GetCloneURL())
	fmt.Println("clone URL:", repo.GetCloneURL())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not clone '%s': %v", name, err)
	}

	fmt.Printf("'%s' cloned successfully\n", name)
	return nil
}

func add(dir string) (err error) {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git add: %v", err)
	}
	return nil
}

func commit(dir string, commitMessage string) (err error) {
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git commit: %v", err)
	}
	return nil
}

func push(dir string) (err error) {
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git push: %v", err)
	}
	return nil
}

func copyFiles(target string, files ...string) (err error) {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("could not copy file '%s' to '%s': %v", file, target, err)
		}

		if err = os.WriteFile(filepath.Join(target, file), data, 0644); err != nil {
			return fmt.Errorf("could not copy file '%s' to '%s': %v", file, target, err)
		}
	}
	return nil
}

// Copies `files` into `repoName` and pushes them to the remote. Clones the repo if necessary.
func UploadFiles(repoName string, commitMessage string, files ...string) (err error) {
	if err := Clone(repoName); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	if err = copyFiles(repoName, files...); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	if err = add(repoName); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	if err = commit(repoName, commitMessage); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	if err = push(repoName); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	return nil
}
