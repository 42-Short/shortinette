package main

import (
	"fmt"
	"net/http"
	"errors"

	internalErrors "github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/pkg/git"
	"github.com/42-Short/shortinette/pkg/tester"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type TeamRequest struct {
	TeamName      string   `json:"teamName"`
	Collaborators []string `json:"collaborators"`
}

func createNewTeam(c *gin.Context) {
	var request TeamRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	if err := git.Create("shortinette-test"); err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Failed to create team"})
	}

	for _, collaborator := range request.Collaborators {
		if err := git.AddCollaborator(request.TeamName, collaborator, "push"); err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": "Failed to add collaborator: " + collaborator})
			return
		}
	}
	c.JSON(200, gin.H{"message": "Team and collaborators created successfully"})
}

func testSubmission(c *gin.Context) {
	repoId := c.Param("repoId")

	if err := tester.Run("testconfig/R00.yaml", repoId, "studentcode"); err != nil {
		var submissionErr *internalErrors.SubmissionError
		if errors.As(err, &submissionErr) {
			switch {
			case errors.Is(submissionErr.Err, internalErrors.ErrEmptyRepo):
				c.JSON(http.StatusBadRequest, gin.H{"error": "The repository is empty"})
			case errors.Is(submissionErr.Err, internalErrors.ErrForbiddenItem):
				c.JSON(http.StatusBadRequest, gin.H{"error": "Forbidden items used in the repository"})
			case errors.Is(submissionErr.Err, internalErrors.ErrInvalidOutput):
				c.JSON(http.StatusBadRequest, gin.H{"error": "The output of the code is invalid"})
			case errors.Is(submissionErr.Err, internalErrors.ErrInvalidCompilation):
				c.JSON(http.StatusBadRequest, gin.H{"error": "The code could not compile"})
			case errors.Is(submissionErr.Err, internalErrors.ErrRuntime):
				c.JSON(http.StatusInternalServerError, gin.H{"error": "The code did not execute as expected"})
			case errors.Is(submissionErr.Err, internalErrors.ErrFailedTests):
				c.JSON(http.StatusBadRequest, gin.H{"error": "The tests failed"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "An unknown error occurred"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unknown error occurred"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tests run successfully"})
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}
	router := gin.Default()

	router.POST("/teams/new", createNewTeam)
	router.GET("/test/:repoId", testSubmission)

	router.Run("0.0.0.0:8080")
}
