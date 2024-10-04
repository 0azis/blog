package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type post struct {
	db *sqlx.DB
}

func (p post) Create(post domain.Post) error {
	_, err := p.db.Query(`insert into posts (user_id, category_id, content) values (?, ?, ?, ?)`, post.UserID, post.CategoryID, post.Content)
	return err
}
