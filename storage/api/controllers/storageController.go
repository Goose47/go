package controllers

import (
	"Goose47/storage/api/errs"
	"Goose47/storage/models"
	"Goose47/storage/utils/repositories"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mime/multipart"
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

type SetForm struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Ttl  int                   `form:"ttl" binding:"required"`
}

func (*StorageController) Set(c *gin.Context) {
	key := c.Param("key")

	var form SetForm
	if err := c.ShouldBindWith(&form, binding.Form); err != nil {
		c.JSON(422, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	id, err := storageRepository.Set(key, &models.StorageItem{
		Path:         "path",
		OriginalName: "orname",
		Ttl:          form.Ttl,
	})

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
