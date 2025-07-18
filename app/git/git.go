// `git` is the package responsible for interactions with GitHub. The repos are _always_ assumed to be in the
// organisation specified by the `GITHUB_ORGANISATION` environment variable.
package git

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/42-Short/shortinette/logger"
	"github.com/google/go-github/v66/github"
)

type GithubService struct {
	Client   *github.Client
	Orga     string
	Token    string
	BasePath string
}

func NewGithubService(authToken string, orga string, basePath string) *GithubService {
	return &GithubService{
		Client:   github.NewClient(nil).WithAuthToken(authToken),
		Orga:     orga,
		Token:    authToken,
		BasePath: basePath,
	}
}

func (gh *GithubService) deleteRepo(name string) (err error) {
	if resp, err := gh.Client.Repositories.Delete(context.Background(), gh.Orga, name); err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("could not delete repo '%s': %v", name, err)
		} else {
			logger.Warning.Printf("repo '%s' not found in orga '%s'\n", name, gh.Orga)
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
			if strings.Contains(strings.ToLower(e.Message), "name already exists on this account") {
				return true
			}
		}
	}
	return false
}

// Creates a new repository `name` under the GitHub organisation.
// If `private` is true, the repository's visibility will be private.
func (gh *GithubService) NewRepo(templateRepoName string, name string, private bool, description string) (err error) {
	includeAllBranches := true
	createdRepo, response, err := gh.Client.Repositories.CreateFromTemplate(context.Background(), gh.Orga, templateRepoName, &github.TemplateRepoRequest{Name: &name, Private: &private, Owner: &gh.Orga, Description: &description, IncludeAllBranches: &includeAllBranches})
	if err != nil {
		if response != nil && response.StatusCode == http.StatusUnprocessableEntity {
			if isRepoAlreadyExists(err) {
				logger.Info.Printf("repo %s already exists under orga %s, skipping\n", name, gh.Orga)
				return nil
			}
		}
		return fmt.Errorf("could not create repo %s: %v", name, err)
	}

	logger.Info.Printf("repo created: %s at URL: %s\n", *createdRepo.Name, *createdRepo.HTMLURL)

	// if err := gh.Clone(name); err != nil {
	// 	logger.Error.Printf("could not clone '%s': %v - 'traces' branch not created", name, err)
	// 	return nil
	// } else {
	// 	if err := gh.UploadFiles(name, "Create 'traces' branch", "traces", true); err != nil {
	// 		logger.Error.Printf("could not create 'traces' branch on repo '%s': %v", name, err)
	// 		return nil
	// 	}
	// }

	// logger.Info.Printf("branch 'traces' successfully created on repo '%s'", name)

	return nil
}

// Adds collaborator `collaboratorName` to repo `repoName` (under the GitHub organisation)
// with access level `permission“.
func (gh *GithubService) AddCollaborator(repoName string, collaboratorName string, permission string) (err error) {
	options := &github.RepositoryAddCollaboratorOptions{
		Permission: permission,
	}

	logger.Info.Printf("adding collaborator %s to repo %s\n", collaboratorName, repoName)

	if _, _, err = gh.Client.Repositories.AddCollaborator(context.Background(), gh.Orga, repoName, collaboratorName, options); err != nil {
		return fmt.Errorf("could not add collaborator %s to repo %s: %v", collaboratorName, repoName, err)
	}

	logger.Info.Printf("user %s added to repo %s with %s access\n", collaboratorName, repoName, permission)
	return nil
}

// Clones repo `name` (from the GitHub organisation).
// Does nothing if the directory is cloned already.
func (gh *GithubService) Clone(name string) (err error) {
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		logger.Info.Printf("'%s' seems to be cloned already, returning\n", name)
		return nil
	}

	cloneURL := fmt.Sprintf("https://%s@github.com/%s/%s.git", gh.Token, gh.Orga, name)

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

func copyDirectory(src string, dest string) (err error) {
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relativePath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})

	return err
}

func copyFile(src string, dest string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err == nil {
		err = os.Chmod(dest, srcInfo.Mode())
	}

	return err
}

func copyFiles(target string, files ...string) (err error) {
	for _, file := range files {

		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("could not stat '%s': %v", file, err)
		}

		if info.IsDir() {
			err = copyDirectory(file, filepath.Join(target, filepath.Base(file)))
			if err != nil {
				return fmt.Errorf("could not copy directory '%s': %v", file, err)
			}
		} else {
			err = copyFile(file, filepath.Join(target, filepath.Base(file)))
			if err != nil {
				return fmt.Errorf("could not copy file '%s' to '%s': %v", file, target, err)
			}
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
		return fmt.Errorf("git checkout %s: %s", to, err)
	}

	cmd = exec.Command("git", "push", "--set-upstream", "origin", to)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git push --set-upstream origin %s: %v", to, err)
	}

	return nil
}

func (gh *GithubService) NewBranch(repoName string, branch string) (err error) {
	if err := gh.Clone(repoName); err != nil {
		return fmt.Errorf("could not create branch '%s' on repo '%s': %v", branch, repoName, err)
	}
	if err := checkout(repoName, branch, true); err != nil {
		return fmt.Errorf("could not create branch '%s' on repo '%s': %v", branch, repoName, err)
	}
	return nil
}

// Copies `files` into `repoName` and pushes them to the branch `branchName` on the remote.
//
// If `createBranch` is set to true, a new branch will be created.
//
// Clones the repo if necessary.
func (gh *GithubService) UploadFiles(repoName string, commitMessage string, branch string, createBranch bool, files ...string) (err error) {
	if err := gh.Clone(repoName); err != nil {
		return fmt.Errorf("could not upload files to '%s': %v", repoName, err)
	}

	defer func() {
		if err := os.RemoveAll(repoName); err != nil {
			logger.Error.Printf("could not tear down repo '%s': %v", repoName, err)
		}
	}()

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
func (gh *GithubService) NewRelease(repoName string, tagName string, releaseName string, body string) (err error) {
	makeLatest := "true"
	if _, _, err := gh.Client.Repositories.CreateRelease(context.Background(), gh.Orga, repoName, &github.RepositoryRelease{
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

// DoesAccountExist checks if the provided GitHub username exists.
// returns a bool indicating if the Account exists
//
// WARNING: Returns an error when the github api request was not successful
func DoesAccountExist(username string) (bool, error) {
	client := github.NewClient(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, response, err := client.Users.Get(ctx, username)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("github API error for username '%s': %v", username, err)
	}
	return *user.Type == "User", nil
}

func (gh *GithubService) CreateModuleTemplate(module int) (templateName string, err error) {
	isTemplate := true
	templateName = fmt.Sprintf("module-0%d-template", module)

	_, response, err := gh.Client.Repositories.Create(context.Background(), gh.Orga, &github.Repository{Name: &templateName, IsTemplate: &isTemplate})
	if err != nil {
		if response != nil && response.StatusCode == http.StatusUnprocessableEntity {
			if isRepoAlreadyExists(err) {
				logger.Info.Printf("repo %s already exists under orga %s, skipping\n", templateName, gh.Orga)
				return templateName, nil
			}
		}
		return "", fmt.Errorf("could not create repo %s: %v", templateName, err)
	}

	if err = gh.Clone(templateName); err != nil {
		return "", fmt.Errorf("could not clone template repo '%s': %v", templateName, err)
	}
	defer func() {
		if err := os.RemoveAll(templateName); err != nil {
			logger.Warning.Printf("could not clean up directory '%s'", templateName)
		}
	}()

	subjectPath := filepath.Join("rust", "subjects", fmt.Sprintf("0%d", module), "README.md")
	devcontainerConfigPath := filepath.Join("rust", ".devcontainer")

	if err = gh.UploadFiles(templateName, fmt.Sprintf("add: devcontainer config + subject for module 0%d", module), "main", false, subjectPath, devcontainerConfigPath); err != nil {
		_ = gh.deleteRepo(templateName)
		return "", fmt.Errorf("could not upload files: %v", err)
	}

	if err = gh.NewBranch(templateName, "traces"); err != nil {
		_ = gh.deleteRepo(templateName)
		return "", fmt.Errorf("could not create 'traces' branch: %v", err)
	}

	return templateName, nil
}
