package repository

import (
	"github.com/keremkartal/goticketra/internal/identity/domain"
	"gorm.io/gorm"
)

// UserRepository arayüzü (Interface): Bağımlılığı soyutlamak için.
// İleride Mock test yazarken işimize çok yarayacak.
type UserRepository interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

// userRepository struct'ı, interface'i implemente eder.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository: Yeni bir repository örneği oluşturur.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	// First: İlk bulduğu kaydı getirir. Bulamazsa hata döner.
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}