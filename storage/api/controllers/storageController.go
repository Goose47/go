package controllers

import "github.com/gin-gonic/gin"

type StorageController struct{}

func (*StorageController) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get",
	})
}

func (*StorageController) Set(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Set",
	})
}

func (*StorageController) Delete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Delete",
	})
}
