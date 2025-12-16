package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/keremkartal/goticketra/internal/identity/domain"
	"github.com/keremkartal/goticketra/internal/identity/repository"
	"github.com/keremkartal/goticketra/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password string) (string, error)
	Login(email, password string) (string, error)
}

type authService struct {
	repo   repository.UserRepository
	config config.Config
}

func NewAuthService(repo repository.UserRepository, cfg config.Config) AuthService {
	return &authService{repo: repo, config: cfg}
}

// Register: Kullanıcıyı kaydeder
func (s *authService) Register(email, password string) (string, error) {
	// 1. Şifreyi Hashle (Güvenlik için asla düz metin saklanmaz)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &domain.User{
		Email:    email,
		Password: string(hashedPass),
	}

	// 2. Repository'e gönder
	if err := s.repo.CreateUser(newUser); err != nil {
		return "", errors.New("kullanıcı oluşturulamadı (email kullanılıyor olabilir)")
	}

	return "Kullanıcı başarıyla oluşturuldu", nil
}

// Login: Giriş yapar ve JWT döner
func (s *authService) Login(email, password string) (string, error) {
	// 1. Kullanıcıyı bul
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("kullanıcı bulunamadı")
	}

	// 2. Şifreyi kontrol et (Hashlenmiş şifre ile girilen şifreyi kıyasla)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("hatalı şifre")
	}

	// 3. JWT Token Üret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.config.JWTExpiration)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}