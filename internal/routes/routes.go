package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	Auth *handlers.AuthHandler
}

func Setup(router *gin.Engine, h *Handlers) {
	v1 := router.Group("/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	registerAuthRoutes(v1, h.Auth)
}
