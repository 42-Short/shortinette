package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/dao"
	"github.com/42-Short/shortinette/logger"
	"github.com/42-Short/shortinette/short"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func launchShort(participantDao *dao.DAO[dao.Participant], config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		participants, err := participantDao.GetAll(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not fetch participants: %v", err)})
			return
		}

		sh := short.NewShort(participants, config)

		if err := sh.Launch(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not launch Short: %v", err)})
			return
		}
	}
}

func githubWebhookHandler(moduleDao *dao.DAO[dao.Module], participantDao *dao.DAO[dao.Participant], config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		var payload gitHubWebhookPayload
		if err := c.ShouldBindBodyWith(&payload, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to bind JSON: %v", err)})
			return
		}
		logger.Info.Printf("got webhook payload from %s on repo %s", payload.Pusher.Name, payload.Repository.Name)

		err := processGithubPayload(payload, moduleDao, participantDao, config)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, payload)
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

		mg := newModuleGrader(moduleDao, participantDao, context.TODO(), config)
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
		logger.Info.Printf("invalid payload (not on main), payload.Ref: %s\n", payload.Ref)
		return nil
	}

	if payload.Commit.Message != "grademe" {
		logger.Info.Printf("invalid payload (commit msg not grademe)\n")
		return nil
	}

	if len(payload.Repository.Name) < len(payload.Pusher.Name) {
		logger.Info.Printf("invalid payload (weird repo name)\n")
		return fmt.Errorf("invalid Repository name: %s", payload.Repository.Name)
	}

	moduleId, err := strconv.Atoi(payload.Repository.Name[len(payload.Repository.Name)-2:])
	if err != nil {
		logger.Info.Printf("invalid payload (broken repo name, no int in the end)\n")
		return fmt.Errorf("invalid Repository name: %s", payload.Repository.Name)
	}

	logger.Info.Printf("push event on %s identified as submission.", payload.Repository.Name)
	mg := newModuleGrader(moduleDao, participantDao, context.TODO(), config)
	go func() {
		err := mg.process(payload.Repository.Name[:len(payload.Repository.Name)-3], moduleId)
		if err != nil {
			logger.Error.Printf("grading failed for %s-%02d: %v", payload.Pusher.Name, moduleId, err)
		}
	}()
	return nil
}
