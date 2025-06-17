package article

import (
	"News_site/internal/models"
	"errors"
)

type Service interface {
	GetByID(id uint16) (*models.Article, error)
	GetAll() ([]models.Article, error)
	Create(title, anons, fullText string) error
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

func (s *service) Create(title, anons, fullText string) error {
	// Можно добавить бизнес-логику перед созданием
	if title == "" || anons == "" || fullText == "" {
		return errors.New("not all data has been filled in")
	}

	article := models.Article{
		Title:    title,
		Anons:    anons,
		FullText: fullText,
	}

	return s.repo.CreateArticle(&article)
}
