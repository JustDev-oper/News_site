package core

import (
	"News_site/internal/services/article"
	"database/sql"
	"log"
)

type Container struct {
	DB             *sql.DB
	ArticleService article.Service
}

func NewContainer(db *sql.DB) *Container {
	// Инициализация всех репозиториев
	articleRepo := article.NewRepository(db)

	// Инициализация сервисов с зависимостями
	return &Container{
		DB:             db,
		ArticleService: article.NewService(articleRepo),
	}
}

// Close освобождает ресурсы
func (c *Container) Close() {
	if err := c.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
