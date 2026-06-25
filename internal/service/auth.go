package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stray4x/ecom-crm/internal/config"
	"github.com/stray4x/ecom-crm/internal/dto"
	"github.com/stray4x/ecom-crm/internal/models"
	"github.com/stray4x/ecom-crm/internal/redis"
	"github.com/stray4x/ecom-crm/internal/repository"
	"github.com/stray4x/ecom-crm/pkg/csrf"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
	CSRFToken    string
}

type AuthService interface {
	Register(req dto.RegisterCustomerRequest) (*dto.CustomerResponse, *Tokens, error)
	Login(req dto.LoginRequest) (*dto.CustomerResponse, *Tokens, error)
	Refresh(refreshToken string) (*Tokens, error)
	Logout(refreshToken string) error
}

type authService struct {
	customerRepo     repository.CustomerRepository
	jwtAccessSecret  string
	jwtRefreshSecret string
	tokenStore       redis.TokenStore
}

func NewAuthService(
	customerRepo repository.CustomerRepository,
	cfg *config.Config,
	tokenStore redis.TokenStore,
) AuthService {
	return &authService{
		customerRepo:     customerRepo,
		tokenStore:       tokenStore,
		jwtAccessSecret:  cfg.JWTAccessSecret,
		jwtRefreshSecret: cfg.JWTRefreshSecret,
	}
}

func (s *authService) Register(req dto.RegisterCustomerRequest) (*dto.CustomerResponse, *Tokens, error) {
	existing, err := s.customerRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, err
	}

	if existing != nil {
		return nil, nil, errors.New("Email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	customer := &models.Customer{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hash),
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(customer.ID)

	if err != nil {
		return nil, nil, err
	}

	return toCustomerResponse(customer), tokens, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.CustomerResponse, *Tokens, error) {
	customer, err := s.customerRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New("Email or password does not match")
	}
	if customer == nil {
		return nil, nil, errors.New("Email or password does not match")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("Email or password does not match")
	}

	tokens, err := s.generateTokens(customer.ID)

	if err != nil {
		return nil, nil, err
	}

	return toCustomerResponse(customer), tokens, nil
}

func (s *authService) Refresh(refreshToken string) (*Tokens, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid refresh token")
		}
		return []byte(s.jwtRefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	customerID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	stored, err := s.tokenStore.Get(context.Background(), customerID.String())
	if err != nil || stored != refreshToken {
		return nil, errors.New("refresh token revoked")
	}

	_ = s.tokenStore.Delete(context.Background(), customerID.String())

	return s.generateTokens(customerID)
}

func (s *authService) Logout(refreshToken string) error {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtRefreshSecret), nil
	})
	if err != nil || !token.Valid {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	if claims["type"] != "refresh" {
		return nil
	}

	customerID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil
	}

	return s.tokenStore.Delete(context.Background(), customerID.String())
}

func (s *authService) generateAccessToken(customerID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":  customerID.String(),
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtAccessSecret))
}

func (s *authService) generateRefreshToken(customerID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":  customerID.String(),
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtRefreshSecret))
}

func (s *authService) generateTokens(customerID uuid.UUID) (*Tokens, error) {
	accessToken, err := s.generateAccessToken(customerID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(customerID)
	if err != nil {
		return nil, err
	}

	if err := s.tokenStore.Save(
		context.Background(),
		customerID.String(),
		refreshToken,
		7*24*time.Hour,
	); err != nil {
		return nil, err
	}

	csrfToken, err := csrf.GenerateToken()
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CSRFToken:    csrfToken,
	}, nil
}

func toCustomerResponse(c *models.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:        c.ID.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
	}
}
