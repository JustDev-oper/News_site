package user

import (
	"News_site/internal/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUserByEmail(email string) (*models.User, error)
	ValidatePassword(user *models.User, password string) bool
	CreateUser(email, password, username string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *service) ValidatePassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *service) CreateUser(email, password, username string) error {
	if email == "" || password == "" || username == "" {
		return errors.New("email or password is empty")
	}
	user := models.RegisterRequest{
		Email:    email,
		Password: password,
		Username: username,
	}

	return s.repo.CreateUser(&user)
}
