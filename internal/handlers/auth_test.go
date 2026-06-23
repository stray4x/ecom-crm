package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stray4x/ecom-crm/internal/config"
	"github.com/stray4x/ecom-crm/internal/service/mocks"
)

func setupTestRouter(h *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(gin.Recovery())

	auth := r.Group("/v1/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/token/refresh", h.RefreshToken)
	}

	return r
}

func TestAuthHandler_Login_Unauthorized(t *testing.T) {
	mockService := new(mocks.MockAuthService)

	handler := NewAuthHandler(mockService, &config.Config{
		AppDomain: "localhost",
		AppEnv:    "test",
	})

	router := setupTestRouter(handler)

	reqBody := `{
		"email":"john42@example.com",
		"password":"wrong-password"
	}`

	mockService.On("Login", mock.Anything).
		Return(nil, nil, errors.New("Email or password does not match"))

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Email or password does not match")

	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_MissingCookie(t *testing.T) {
	mockService := new(mocks.MockAuthService)

	handler := NewAuthHandler(mockService, &config.Config{
		AppDomain: "localhost",
		AppEnv:    "test",
	})

	router := setupTestRouter(handler)

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/token/refresh", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "refresh token missing")

	mockService.AssertExpectations(t)
}
