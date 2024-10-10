package domain

import "time"

type Post struct {
	ID         int    `json:"ID" db:"id"`
	UserID     int    `json:"userID" db:"user_id"`
	CategoryID int    `json:"categoryID" db:"category_id"`
	Date       string `json:"date" db:"date"`
	Content    string `json:"content" db:"content"`
	Public     bool   `json:"-" db:"public"`
}

type UserPost struct {
	ID           int       `json:"ID" db:"id"`
	Username     string    `json:"username" db:"username"`
	CategoryName string    `json:"categoryName" db:"category_name"`
	Name         *string   `json:"name" db:"name"`
	Date         time.Time `json:"date" db:"date"`
	Content      string    `json:"content" db:"content"`
}

type PostCredentials struct {
	UserID  int    `json:"userID"`
	Content string `json:"content"`
}
