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

// @Summary Register new customer
// @Description Creates new customer account. Server sets cookies (refresh_token, csrf_token)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterCustomerRequest true "Register request"
// @Success 201 {object} dto.AuthResponse
// @Header 200 {string} Set-Cookie "refresh_token httpOnly"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/register [post]
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

// @Summary Login existing customer
// @Description Authenticate user. Server sets cookies (refresh_token, csrf_token).
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthResponse
// @Header 200 {string} Set-Cookie "refresh_token httpOnly"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/login [post]
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

// @Summary Logout customer
// @Description
// Logs out the current customer by invalidating the refresh_token session
// and clearing authentication cookies (refresh_token).
//
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.LogoutResponse
// @Header 200 {string} Set-Cookie "refresh_token httpOnly; Max-Age=0"
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err == nil && refreshToken != "" {
		h.authService.Logout(refreshToken)
	}

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

// @Summary Refresh access token
// @Description
// Issues new access token using refresh_token stored in HttpOnly cookie.
//
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /v1/auth/token/refresh [post]
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
