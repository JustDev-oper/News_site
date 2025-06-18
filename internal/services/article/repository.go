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
	GetArticlesByUserID(userID uint) ([]models.Article, error)
	DeleteArticle(id uint16) error
	UpdateArticle(article *models.Article) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (repo *repository) GetAllArticles() ([]models.Article, error) {
	var articles []models.Article
	query := `
		SELECT a.id, a.title, a.anons, a.full_text, a.user_id, a.created_at,
		       u.id, u.email, u.username, u.created_at
		FROM articles a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
	`

	res, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying articles: %v", err)
	}
	defer res.Close()

	for res.Next() {
		var article models.Article
		var user models.User
		err = res.Scan(
			&article.Id, &article.Title, &article.Anons, &article.FullText,
			&article.UserID, &article.CreatedAt,
			&user.ID, &user.Email, &user.Username, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning article: %v", err)
		}
		article.User = &user
		articles = append(articles, article)
	}

	if err = res.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %v", err)
	}

	return articles, nil
}

func (repo *repository) CreateArticle(article *models.Article) error {
	_, err := repo.db.Exec(
		"INSERT INTO `articles` (`title`, `anons`, `full_text`, `user_id`, `created_at`) VALUES (?, ?, ?, ?, ?)",
		article.Title,
		article.Anons,
		article.FullText,
		article.UserID,
		article.CreatedAt,
	)
	return err
}

func (repo *repository) GetArticleByID(id uint16) (*models.Article, error) {
	var article models.Article
	var user models.User

	query := `
		SELECT a.id, a.title, a.anons, a.full_text, a.user_id, a.created_at,
		       u.id, u.email, u.username, u.created_at
		FROM articles a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.id = ?
	`

	err := repo.db.QueryRow(query, id).Scan(
		&article.Id, &article.Title, &article.Anons, &article.FullText,
		&article.UserID, &article.CreatedAt,
		&user.ID, &user.Email, &user.Username, &user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("article not found")
		}
		return nil, err
	}

	article.User = &user
	return &article, nil
}

func (repo *repository) GetArticlesByUserID(userID uint) ([]models.Article, error) {
	var articles []models.Article
	query := `
		SELECT a.id, a.title, a.anons, a.full_text, a.user_id, a.created_at,
		       u.id, u.email, u.username, u.created_at
		FROM articles a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.user_id = ?
		ORDER BY a.created_at DESC
	`

	res, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user articles: %v", err)
	}
	defer res.Close()

	for res.Next() {
		var article models.Article
		var user models.User
		err = res.Scan(
			&article.Id, &article.Title, &article.Anons, &article.FullText,
			&article.UserID, &article.CreatedAt,
			&user.ID, &user.Email, &user.Username, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user article: %v", err)
		}
		article.User = &user
		articles = append(articles, article)
	}

	if err = res.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user articles: %v", err)
	}

	return articles, nil
}

func (repo *repository) DeleteArticle(id uint16) error {
	_, err := repo.db.Exec(
		"DELETE FROM `articles` WHERE `id` = ?",
		id,
	)
	return err
}

func (repo *repository) UpdateArticle(article *models.Article) error {
	_, err := repo.db.Exec(
		"UPDATE `articles` SET `title` = ?, `anons` = ?, `full_text` = ? WHERE `id` = ? AND `user_id` = ?",
		article.Title,
		article.Anons,
		article.FullText,
		article.Id,
		article.UserID,
	)
	return err
}
