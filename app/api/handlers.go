package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/dao"
	"github.com/42-Short/shortinette/logger"
	"github.com/gin-gonic/gin"
)

type gitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Commit struct {
		Message string `json:"message"`
	} `json:"head_commit"`
}

func githubWebhookHandler(moduleDao *dao.DAO[dao.Module], participantDao *dao.DAO[dao.Participant], config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload gitHubWebhookPayload

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = processGithubPayload(payload, moduleDao, participantDao, config)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusProcessing, payload)
	}
}

func gradingHandler(moduleDao *dao.DAO[dao.Module], participantDao *dao.DAO[dao.Participant], config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		args := collectArgs(c.Params)
		module, err := moduleDao.Get(ctx, args...)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get %s: %v: %v", moduleDao.Name(), args, err)})
			return
		}

		mg := newModuleGrader(moduleDao, participantDao, context.TODO(), config);
		go func() {
            err := mg.process(module.IntraLogin, module.Id)
            if err != nil {
                logger.Error.Printf("grading failed for %s%d: %v", module.IntraLogin, module.Id, err)
            }
        }()
		c.JSON(http.StatusProcessing, fmt.Sprintf("grading %s%d...", module.IntraLogin, module.Id))
	}
}

func insertItemHandler[T any](dao *dao.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item T
		err := c.ShouldBindJSON(&item)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = dao.Insert(ctx, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to insert %s: %v", dao.Name(), err)})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func updateItemHandler[T any](dao *dao.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item T

		err := c.ShouldBindJSON(&item)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = dao.Update(ctx, item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete %s: %v: %v", dao.Name(), item, err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted item %v from %s", item, dao.Name())})
	}
}

func getAllItemsHandler[T any](dao *dao.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		items, err := dao.GetAll(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get all %s`s: %v", dao.Name(), err)})
			return
		}
		if len(items) == 0 {
			c.JSON(http.StatusNoContent, items)
		} else {
			c.JSON(http.StatusOK, items)
		}
	}
}

func getItemHandler[T any](dao *dao.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		args := collectArgs(c.Params)
		item, err := dao.Get(ctx, args...)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to get %s: %v: %v", dao.Name(), args, err)})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func deleteItemHandler[T any](dao *dao.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		args := collectArgs(c.Params)
		err := dao.Delete(ctx, args...)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("failed to delete %s: %v", dao.Name(), err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted item from %s %v", dao.Name(), args)})
	}
}

func collectArgs(params gin.Params) []any {
	args := make([]any, 0, len(params))

	for _, param := range params {
		args = append(args, param.Value)
	}
	return args
}

func processGithubPayload(payload gitHubWebhookPayload, moduleDao *dao.DAO[dao.Module], participantDao *dao.DAO[dao.Participant], config config.Config) error {
	if payload.Ref != "refs/heads/main" || payload.Pusher.Name == os.Getenv("GITHUB_ADMIN") {
		return nil
	}
	if strings.ToLower(payload.Commit.Message) != "grademe" {
		return nil
	}

	if len(payload.Repository.Name) < len(payload.Pusher.Name) {
		return fmt.Errorf("invalid Repository name: %s", payload.Repository.Name)
	}
	moduleId, err := strconv.Atoi(payload.Repository.Name[len(payload.Pusher.Name):])
	if err != nil {
		return fmt.Errorf("invalid Repository name: %s", payload.Repository.Name)
	}

	logger.Info.Printf("push event on %s identified as submission.", payload.Repository.Name)
	mg := newModuleGrader(moduleDao, participantDao, context.TODO(), config);
	go func() {
		err := mg.process(payload.Pusher.Name, moduleId)
		if err != nil {
			logger.Error.Printf("grading failed for %s%d: %v", payload.Pusher.Name, moduleId, err)
		}
	}()
	return nil
}
