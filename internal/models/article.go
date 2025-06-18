package models

import (
	"time"
)

type Article struct {
	Id        uint16    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Anons     string    `json:"anons" db:"anons"`
	FullText  string    `json:"full_text" db:"full_text"`
	UserID    uint      `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	User      *User     `json:"user,omitempty" db:"-"` // Связь с пользователем
	LikesCount int `json:"likes_count" db:"-"`
}
