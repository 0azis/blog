package domain

import (
	"time"
)

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

func (up UserPost) Validate() bool {
	return *up.Title != "" && *up.Content != ""
}

type PostCredentials struct {
	UserID  int    `json:"userID"`
	Title   string `json:"title"`
	Preview string `json:"image"`
	Content string `json:"content"`
}
