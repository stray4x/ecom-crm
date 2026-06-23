package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stray4x/ecom-crm/internal/config"
	"github.com/stray4x/ecom-crm/internal/dto"
	"github.com/stray4x/ecom-crm/internal/models"
	"github.com/stray4x/ecom-crm/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

var cfg_test = &config.Config{
	JWTAccessSecret:  "jwt-access-secret",
	JWTRefreshSecret: "jwt-refresh-secret",
}

func hashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name string

		request dto.RegisterCustomerRequest

		existingUser *models.Customer
		getErr       error

		createErr error

		expectError  bool
		errorMessage string
	}{
		{
			name: "success",
			request: dto.RegisterCustomerRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Phone:     "+380991234567",
				Password:  "password123",
			},
			existingUser: nil,
			getErr:       nil,
			createErr:    nil,
			expectError:  false,
		},
		{
			name: "email already exists",
			request: dto.RegisterCustomerRequest{
				Email: "john@example.com",
			},
			existingUser: &models.Customer{Email: "john@example.com"},
			getErr:       nil,
			expectError:  true,
			errorMessage: "Email already in use",
		},
		{
			name: "repo GetByEmail error",
			request: dto.RegisterCustomerRequest{
				Email: "john@example.com",
			},
			existingUser: nil,
			getErr:       errors.New("db error"),
			expectError:  true,
		},
		{
			name: "repo Create error",
			request: dto.RegisterCustomerRequest{
				Email: "john@example.com",
			},
			existingUser: nil,
			getErr:       nil,
			createErr:    errors.New("insert failed"),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockCustomerRepository(t)

			mockRepo.EXPECT().
				GetByEmail(tt.request.Email).
				Return(tt.existingUser, tt.getErr)

			if tt.existingUser == nil && tt.getErr == nil {
				mockRepo.EXPECT().
					Create(mock.AnythingOfType("*models.Customer")).
					Return(tt.createErr)
			}

			service := NewAuthService(mockRepo, cfg_test)

			res, tokens, err := service.Register(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.Nil(t, tokens)

				if tt.errorMessage != "" {
					assert.Equal(t, tt.errorMessage, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.NotNil(t, tokens)
			assert.NotEmpty(t, tokens.AccessToken)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		password     string
		mockUser     *models.Customer
		mockErr      error
		expectError  bool
		errorMessage string
	}{
		{
			name:     "success",
			email:    "john@example.com",
			password: "password123",
			mockUser: &models.Customer{
				Email:        "john@example.com",
				PasswordHash: hashPassword("password123"),
			},
			expectError: false,
		},
		{
			name:         "user not found",
			email:        "john@example.com",
			password:     "password123",
			mockUser:     nil,
			expectError:  true,
			errorMessage: "Email or password does not match",
		},
		{
			name:     "wrong password",
			email:    "john@example.com",
			password: "wrong",
			mockUser: &models.Customer{
				Email:        "john@example.com",
				PasswordHash: hashPassword("correct"),
			},
			expectError:  true,
			errorMessage: "Email or password does not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockCustomerRepository(t)

			mockRepo.EXPECT().
				GetByEmail(tt.email).
				Return(tt.mockUser, tt.mockErr)

			service := NewAuthService(mockRepo, cfg_test)

			res, tokens, err := service.Login(dto.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
				assert.Nil(t, tokens)
				assert.Equal(t, tt.errorMessage, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, res)
			assert.NotNil(t, tokens)
		})
	}
}

func TestAuthService_Refresh_Success(t *testing.T) {
	mockRepo := mocks.NewMockCustomerRepository(t)

	authService := NewAuthService(mockRepo, cfg_test)

	customerID := uuid.New()

	claims := jwt.MapClaims{
		"sub":  customerID.String(),
		"exp":  time.Now().Add(1 * time.Hour).Unix(),
		"type": "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(cfg_test.JWTRefreshSecret))
	assert.NoError(t, err)

	tokens, err := authService.Refresh(refreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, tokens)

	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.NotEmpty(t, tokens.CSRFToken)
}
