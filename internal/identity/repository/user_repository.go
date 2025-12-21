package repository

import (
	"github.com/keremkartal/goticketra/internal/identity/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}