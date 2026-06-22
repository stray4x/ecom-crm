package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/handlers"
)

type Handlers struct {
	Auth *handlers.AuthHandler
}

func Setup(router *gin.Engine, h *Handlers) {
	v1 := router.Group("/v1")

	registerAuthRoutes(v1, h.Auth)
}
