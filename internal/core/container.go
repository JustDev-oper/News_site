package core

import (
	"News_site/internal/services/article"
	"News_site/internal/services/user"
	"database/sql"
	"log"
)

type Container struct {
	DB             *sql.DB
	ArticleService article.Service
	UserService    user.Service
}

func NewContainer(db *sql.DB) *Container {
	// Инициализация всех репозиториев
	articleRepo := article.NewRepository(db)
	userRepo := user.NewRepository(db)

	// Инициализация сервисов с зависимостями
	return &Container{
		DB:             db,
		ArticleService: article.NewService(articleRepo),
		UserService:    user.NewService(userRepo),
	}
}

// Close освобождает ресурсы
func (c *Container) Close() {
	if err := c.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
