package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stray4x/ecom-crm/internal/config"
	"github.com/stray4x/ecom-crm/internal/dto"
	"github.com/stray4x/ecom-crm/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	config      *config.Config
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, config: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, tokens, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	setAuthCookies(c, tokens.RefreshToken, tokens.CSRFToken, h.config.AppDomain, h.config.AppEnv)

	c.JSON(http.StatusCreated, gin.H{
		"customer":    customer,
		"accessToken": tokens.AccessToken,
		"csrfToken":   tokens.CSRFToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, tokens, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setAuthCookies(c, tokens.RefreshToken, tokens.CSRFToken, h.config.AppDomain, h.config.AppEnv)

	c.JSON(http.StatusOK, gin.H{
		"customer":    customer,
		"accessToken": tokens.AccessToken,
		"csrfToken":   tokens.CSRFToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		h.config.AppDomain,
		h.config.AppEnv == "production",
		true,
	)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	tokens, err := h.authService.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	setAuthCookies(c, tokens.RefreshToken, tokens.CSRFToken, h.config.AppDomain, h.config.AppEnv)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": tokens.AccessToken,
		"csrfToken":   tokens.CSRFToken,
	})
}

func setAuthCookies(c *gin.Context, refreshToken, csrfToken, appDomain, appEnv string) {
	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60,
		"/",
		appDomain,
		appEnv == "production",
		true,
	)

	c.SetCookie(
		"csrf_token",
		csrfToken,
		7*24*60*60,
		"/",
		appDomain,
		appEnv == "production",
		false,
	)
}
