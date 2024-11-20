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
		c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented yet"})
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
		c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented yet"})
	}
}

func DeleteItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented yet"})
	}
}
