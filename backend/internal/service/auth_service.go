package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"dalivim/internal/models"
	"dalivim/internal/repository"
)

type AuthService interface {
	Register(email, password, name, role string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(email, password, name, role string) (*models.User, string, error) {
	// Check if user already exists
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, "", errors.New("user already exists")
	}

	// TODO: Hash password with bcrypt
	user := &models.User{
		Email:    email,
		Password: password, // Should be hashed
		Name:     name,
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token := generateToken()
	return user, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// TODO: Compare hashed password with bcrypt
	if user.Password != password {
		return nil, "", errors.New("invalid credentials")
	}

	token := generateToken()
	return user, token, nil
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
