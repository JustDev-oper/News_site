package article

import (
	"News_site/internal/models"
	"errors"
	"time"
)

type Service interface {
	GetByID(id uint16) (*models.Article, error)
	GetAll() ([]models.Article, error)
	Create(title, anons, fullText string, userID uint) error
	GetByUserID(userID uint) ([]models.Article, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(id uint16) (*models.Article, error) {
	return s.repo.GetArticleByID(id)
}

func (s *service) GetAll() ([]models.Article, error) {
	return s.repo.GetAllArticles()
}

func (s *service) Create(title, anons, fullText string, userID uint) error {
	if title == "" || anons == "" || fullText == "" {
		return errors.New("not all data has been filled in")
	}

	if userID == 0 {
		return errors.New("user ID is required")
	}

	article := models.Article{
		Title:     title,
		Anons:     anons,
		FullText:  fullText,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateArticle(&article)
}

func (s *service) GetByUserID(userID uint) ([]models.Article, error) {
	return s.repo.GetArticlesByUserID(userID)
}

// TODO: Сделать удаление и изменение постов, а так же привязать посты к пользователю, который их создал
