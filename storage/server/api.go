package server

import (
	"Goose47/storage/api/controllers"
	"github.com/gin-gonic/gin"
)

func AddApiGroup(r *gin.Engine) {
	api := r.Group("api")
	{
		v1 := api.Group("v1")
		{
			v1.GET("healthcheck")

			storage := v1.Group("storage")
			{
				c := new(controllers.StorageController)
				storage.GET("/:key", c.Get)
				storage.POST("/:key", c.Set)
				storage.DELETE("/:key", c.Delete)
			}
		}
	}
}
