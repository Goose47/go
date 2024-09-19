package server

import (
	"Goose47/storage/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Init() {
	r := NewRouter()
	gin.SetMode(config.AppConfig.Mode)

	addr := fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port)
	log.Fatal(r.Run(addr))
}

func NewRouter() *gin.Engine {
	r := gin.New()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	AddApiGroup(r)

	return r
}
