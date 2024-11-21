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

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/logger"
	"github.com/google/go-github/v66/github"
)

func deleteRepo(name string) (err error) {
	client := github.NewClient(nil).WithAuthToken(config.C.TOKEN_GITHUB)

	if resp, err := client.Repositories.Delete(context.Background(), config.C.ORGA_GITHUB, name); err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("could not delete repo '%s': %v", name, err)
		} else {
			logger.Warning.Printf("repo '%s' not found in orga '%s'\n", name, config.C.ORGA_GITHUB)
			return nil
		}
	}

	logger.Info.Printf("repo '%s' successfully deleted\n", name)
	return nil
}

// Checks whether `err` is related to the repo already existing.
func isRepoAlreadyExists(err error) (exists bool) {
	if githubErr, ok := err.(*github.ErrorResponse); ok {
		for _, e := range githubErr.Errors {
			if strings.Contains(e.Message, "Name already exists on this account") {
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
	client := github.NewClient(nil).WithAuthToken(config.C.TOKEN_GITHUB)

	createdRepo, response, err := client.Repositories.CreateFromTemplate(context.Background(), config.C.ORGA_GITHUB, config.C.TEMPLATE_REPO, &github.TemplateRepoRequest{Name: &name, Private: &private, Owner: &config.C.ORGA_GITHUB, Description: &description})
	if err != nil {
		if response != nil && response.StatusCode == http.StatusUnprocessableEntity {
			if isRepoAlreadyExists(err) {
				logger.Info.Printf("repo %s already exists under orga %s, skipping\n", name, config.C.ORGA_GITHUB)
				return nil
			}
		}
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	logger.Info.Printf("repo created: %s at URL: %s\n", *createdRepo.Name, *createdRepo.HTMLURL)
	return nil
}

// Adds collaborator `collaboratorName` to repo `repoName` (under the GitHub organisation specified by
// the GITHUB_ORGANISATION environment variable) with access level `permission“.
func AddCollaborator(repoName string, collaboratorName string, permission string) (err error) {
	client := github.NewClient(nil).WithAuthToken(config.C.TOKEN_GITHUB)

	options := &github.RepositoryAddCollaboratorOptions{
		Permission: permission,
	}

	if _, _, err = client.Repositories.AddCollaborator(context.Background(), config.C.ORGA_GITHUB, repoName, collaboratorName, options); err != nil {
		return fmt.Errorf("could not add collaborator %s to repo %s: %v", collaboratorName, repoName, err)
	}

	logger.Info.Printf("user %s added to repo %s with %s access\n", collaboratorName, repoName, permission)
	return nil
}

// Clones repo `name` (from the GitHub organisation specified by the GITHUB_ORGANISATION
// environment variable). Does nothing if the directory is cloned already.
func Clone(name string) (err error) {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		logger.Info.Printf("'%s' seems to cloned already, returning\n", name)
		return nil
	}

	cloneURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", config.C.TOKEN_GITHUB, config.C.ORGA_GITHUB, name)

	cmd := exec.Command("git", "clone", cloneURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not clone '%s': %v", name, err)
	}

	logger.Info.Printf("'%s' cloned successfully\n", name)
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

func checkout(dir string, to string, createBranch bool) (err error) {
	if to == "main" {
		return nil
	}

	var cmd *exec.Cmd
	if createBranch {
		cmd = exec.Command("git", "checkout", "-b", to)
	} else {
		cmd = exec.Command("git", "checkout", to)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git checkout %s: %v", to, err)
	}

	cmd = exec.Command("git", "push", "--set-upstream", "origin", to)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git --set-upstream origin %s: %v", to, err)
	}

	return nil
}

// Copies `files` into `repoName` and pushes them to the branch `branchName` on the remote.
//
// If `createBranch` is set to true, a new branch will be created.
//
// Clones the repo if necessary.
func UploadFiles(repoName string, commitMessage string, branch string, createBranch bool, files ...string) (err error) {
	if err := Clone(repoName); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	if err = checkout(repoName, branch, createBranch); err != nil {
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

// Adds a release to `repoName` named `releaseName`, tagged `tagName`, with `body` as body, and
// makes it the latest release for the repo.
//
// WARNING: Tag names must be unique. Append a unique ID (like number of grading attempts) to `releaseName`.
//
// WARNING 2: Returns an error when the repository is empty (due to tarball creation not being possible without
// some content). This should not be an issue for shortinette though, since we always upload subjects when creating
// the repos.
func NewRelease(repoName string, tagName string, releaseName string, body string) (err error) {
	client := github.NewClient(nil).WithAuthToken(config.C.TOKEN_GITHUB)

	makeLatest := "true"
	if _, _, err := client.Repositories.CreateRelease(context.Background(), config.C.ORGA_GITHUB, repoName, &github.RepositoryRelease{
		Name:       &releaseName,
		Body:       &body,
		TagName:    &tagName,
		MakeLatest: &makeLatest,
	}); err != nil {
		return fmt.Errorf("could not add release '%s', tagged '%s' to repo '%s': %v", releaseName, tagName, repoName, err)
	}

	logger.Info.Printf("added release '%s', tagged '%s' to repo '%s'", releaseName, tagName, repoName)
	return nil
}
