package controllers

import (
	"Goose47/storage/api/errs"
	"Goose47/storage/config"
	"Goose47/storage/models"
	"Goose47/storage/utils"
	"Goose47/storage/utils/repositories"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mime/multipart"
	"net/http"
	"path"
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
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	item := &models.StorageItem{}

	item.OriginalName = form.File.Filename
	item.Ttl = form.Ttl
	item.Path = utils.GenerateRandomString(20) + path.Ext(form.File.Filename)

	err := c.SaveUploadedFile(form.File, path.Join(config.FSConfig.Base, item.Path))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
