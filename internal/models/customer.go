package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName    string    `gorm:"not null"`
	LastName     string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	Phone        string    `gorm:"unique"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
