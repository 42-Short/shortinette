package main

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/git"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	if err := git.NewRepo("repo", true, "this is a new repo"); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	defer func() {
		if err := os.RemoveAll("repo"); err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}()

	if err := git.AddCollaborator("repo", "winstonallo", "write"); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if err := git.UploadFiles("repo", "add stuff", "Dockerfile", "go.mod", "go.sum"); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

// func helloWorld(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
// }

// func router() (router *gin.Engine) {
// 	router = gin.Default()

// 	router.GET("/admin", helloWorld)

// 	return router
// }

// func main() {
// 	r := router()
// 	if err := r.Run("0.0.0.0:5000"); err != nil {
// 		fmt.Printf("error running gin server: %v", err)
// 	}
// }
