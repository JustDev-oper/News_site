package models

type Like struct {
	UserID    uint   `json:"user_id" db:"user_id"`
	ArticleID uint16 `json:"article_id" db:"article_id"`
}