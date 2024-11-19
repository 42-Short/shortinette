package server

import (
	"net/http"

	"github.com/42-Short/shortinette/data"
	"github.com/gin-gonic/gin"
)

func InsertItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	}
}

func UpdateItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	}
}

func GetAllItemsHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	}
}

func GetItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	}
}

func DeleteItemHandler[T any](dao *data.DAO[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	}
}
