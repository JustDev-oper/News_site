package article

import (
	"News_site/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type Repository interface {
	GetArticleByID(id uint16) (*models.Article, error)
	GetAllArticles() ([]models.Article, error)
	CreateArticle(article *models.Article) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (repo *repository) GetAllArticles() ([]models.Article, error) {
	var articles []models.Article
	res, err := repo.db.Query("SELECT * FROM `articles`")
	if err != nil {
		return nil, fmt.Errorf("error querying articles: %v", err)
	}
	defer res.Close()

	for res.Next() {
		var article models.Article
		err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil {
			return nil, fmt.Errorf("error scanning article: %v", err)
		}
		articles = append(articles, article)
	}

	if err = res.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %v", err)
	}

	return articles, nil
}

func (repo *repository) CreateArticle(article *models.Article) error {
	_, err := repo.db.Exec(
		"INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES (?, ?, ?)",
		article.Title,
		article.Anons,
		article.FullText,
	)
	return err
}

func (repo *repository) GetArticleByID(id uint16) (*models.Article, error) {
	var article models.Article
	err := repo.db.QueryRow(
		"SELECT id, title, anons, full_text FROM `articles` WHERE `id` = ?",
		id,
	).Scan(
		&article.Id,
		&article.Title,
		&article.Anons,
		&article.FullText,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("article not found")
		}
		return nil, err
	}

	return &article, nil
}
