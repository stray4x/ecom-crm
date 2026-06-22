package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stray4x/ecom-crm/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	GetByID(id uuid.UUID) (*models.Customer, error)
	GetByEmail(email string) (*models.Customer, error)
	Update(customer *models.Customer) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("id = ?", id).First(&customer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &customer, err
}

func (r *customerRepository) GetByEmail(email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("email = ?", email).First(&customer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}
