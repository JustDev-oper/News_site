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
	Delete(id uint16) error
	Update(id uint16, title, anons, fullText string, userID uint) error
	LikeArticle(userID uint, articleID uint16) error
	UnlikeArticle(userID uint, articleID uint16) error
	IsArticleLikedByUser(userID uint, articleID uint16) (bool, error)
	GetLikesCount(articleID uint16) (int, error)
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

func (s *service) GetByUserID(userID uint) ([]models.Article, error) {
	return s.repo.GetArticlesByUserID(userID)
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

func (s *service) Delete(id uint16) error {
	return s.repo.DeleteArticle(id)
}

func (s *service) Update(id uint16, title, anons, fullText string, userID uint) error {
	if title == "" || anons == "" || fullText == "" {
		return errors.New("not all data has been filled in")
	}

	if userID == 0 {
		return errors.New("user ID is required")
	}

	article, err := s.GetByID(id)
	if err != nil {
		return err
	}
	
	if article.UserID != userID {
		return errors.New("access denied: article does not belong to user")
	}

	article.Title = title
	article.Anons = anons
	article.FullText = fullText
	article.UserID = userID
	return s.repo.UpdateArticle(article)
}

func (s *service) LikeArticle(userID uint, articleID uint16) error {
	return s.repo.LikeArticle(userID, articleID)
}

func (s *service) UnlikeArticle(userID uint, articleID uint16) error {
	return s.repo.UnlikeArticle(userID, articleID)
}

func (s *service) IsArticleLikedByUser(userID uint, articleID uint16) (bool, error) {
	return s.repo.IsArticleLikedByUser(userID, articleID)
}

func (s *service) GetLikesCount(articleID uint16) (int, error) {
	return s.repo.GetLikesCount(articleID)
}
