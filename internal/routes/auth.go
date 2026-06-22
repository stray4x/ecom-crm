package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/handlers"
	"github.com/stray4x/ecom-crm/internal/middleware"
)

func registerAuthRoutes(rg *gin.RouterGroup, h *handlers.AuthHandler) {
	auth := rg.Group("/auth")

	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/logout", h.Logout)
	auth.POST("/token/refresh", middleware.CSRFMiddleware(), h.RefreshToken)
}
