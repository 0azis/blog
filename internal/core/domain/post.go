package domain

import (
	"time"
)

type Post struct {
	ID         int    `json:"ID" db:"id"`
	UserID     int    `json:"userID" db:"user_id"`
	CategoryID int    `json:"categoryID" db:"category_id"`
	Title      string `json:"title" db:"title"`
	Preview    string `json:"preview" db:"preview"`
	Date       string `json:"date" db:"date"`
	Content    string `json:"content" db:"content"`
	Public     bool   `json:"-" db:"public"`
}

type UserPost struct {
	ID       int       `json:"id" db:"id"`
	Title    *string   `json:"title" db:"title"`
	Date     time.Time `json:"createdAt" db:"date"`
	Preview  *string   `json:"preview" db:"preview"`
	Username string    `json:"username" db:"username"`
	Name     *string   `json:"name" db:"name"`
	Avatar   *string   `json:"avatar" db:"avatar"`
	Content  *string   `json:"content" db:"content"`
}

type PostCredentials struct {
	UserID  int    `json:"userID"`
	Title   string `json:"title"`
	Preview string `json:"image"`
	Content string `json:"content"`
}
