package dto

type RegisterCustomerRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Password  string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Customer    CustomerResponse `json:"customer"`
	AccessToken string           `json:"accessToken"`
	CSRFToken   string           `json:"csrfToken"`
}

type CustomerResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	CSRFToken   string `json:"csrfToken"`
}
