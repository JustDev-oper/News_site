package models

type Article struct {
	Id       uint16 `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	Anons    string `json:"anons" db:"anons"`
	FullText string `json:"full_text" db:"full_text"`
}
