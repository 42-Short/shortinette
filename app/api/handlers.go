package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/42-Short/shortinette/data"
	"github.com/gin-gonic/gin"
)

func InsertItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item T
		err := c.ShouldBindJSON(&item)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = dao.Insert(context.TODO(), item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to insert %s: %v", dao.Name(), err)})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func UpdateItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var item T

		err := c.ShouldBindJSON(&item)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = dao.Update(context.TODO(), item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete %s: %v", dao.Name(), err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("deleted item from %s %v", dao.Name(), args)})
	}
}

func GetAllItemsHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := dao.GetAll(context.TODO())
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

func GetItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		args := collectArgs(c.Params)
		item, err := dao.Get(context.TODO(), args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get %s %v: %v", dao.Name(), args, err)})
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func DeleteItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		args := collectArgs(c.Params)
		err := dao.Delete(context.TODO(), args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to delete %s: %v", dao.Name(), err)})
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
