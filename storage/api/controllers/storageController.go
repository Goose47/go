package controllers

import (
	"Goose47/storage/api/errs"
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
	"os"
	"path"
	"time"
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

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+item.OriginalName)
	c.File(item.GetFullPath())
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

	// if exp == 0, document never expires
	exp := form.Ttl
	if exp > 0 {
		exp += int(time.Now().Unix())
	}

	item.OriginalName = form.File.Filename
	item.Exp = exp
	item.Path = utils.GenerateRandomString(20) + path.Ext(form.File.Filename)

	err := c.SaveUploadedFile(form.File, item.GetFullPath())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, err := storageRepository.Set(key, item)

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

	if _, err = os.Stat(item.GetFullPath()); err == nil {
		err = os.Remove(item.GetFullPath())
	}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}
