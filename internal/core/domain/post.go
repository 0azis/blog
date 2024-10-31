package domain

import (
	"time"
)

type Post struct {
	ID      int       `json:"id"`
	Title   *string   `json:"title" db:"title"`
	Date    time.Time `json:"createdAt" db:"date"`
	Preview *string   `json:"preview" db:"preview"`
	Author  Author    `json:"author"`
	Content *string   `json:"content" db:"content"`
	Views   int       `json:"views" db:"views"`
	Tags    Tags      `json:"tags"`
}

func (p Post) Validate() bool {
	return *p.Title != "" && *p.Content != ""
}

type PostCredentials struct {
	UserID  int    `json:"userID"`
	Title   string `json:"title"`
	Preview string `json:"image"`
	Content string `json:"content"`
}
