package user

import (
	"News_site/internal/models"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.RegisterRequest) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (repo *repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := repo.db.QueryRow(
		"SELECT id, email, password, username, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (repo *repository) CreateUser(user *models.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(
		"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
		user.Email, string(hashedPassword), user.Username,
	)
	return err
}
