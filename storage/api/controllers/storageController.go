package controllers

import (
	"Goose47/storage/api/errs"
	"Goose47/storage/models"
	"Goose47/storage/utils/repositories"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type StorageController struct{}

var storageRepository *repositories.StorageRepository

func init() {
	storageRepository = new(repositories.StorageRepository)
}

func (*StorageController) Get(c *gin.Context) {
	key := c.Param("key")

	item, err := storageRepository.FindByKey(key)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Error(&errs.NotFoundError{Message: fmt.Sprintf("%s is not found", key)})
			return
		}
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Got %v", item),
	})
}

func (*StorageController) Set(c *gin.Context) {
	key := c.Param("key")

	item := &models.StorageItem{
		Key:          key,
		Path:         "path",
		Ttl:          100,
		OriginalName: "naem",
	}

	id, err := storageRepository.Set(item)

	if err != nil {
		log.Panic(err)
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Set %s", id),
	})
}

func (*StorageController) Delete(c *gin.Context) {
	key := c.Param("key")

	item, err := storageRepository.DeleteByKey(key)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Error(&errs.NotFoundError{Message: fmt.Sprintf("%s is not found", key)})
			return
		}
		log.Panic(err)
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Deleted %v", item),
	})
}
