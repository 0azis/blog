package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type comment struct {
	db *sqlx.DB
}

func (c comment) Create(comment domain.CommentCredentials) error {
	_, err := c.db.Query(`insert into comments (post_id, user_id, text) values (?, ?, ?)`, comment.PostID, comment.UserID, comment.Text)
	return err
}

func (c comment) GetByPostID(postID int) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := c.db.Select(&comments, `select * from comments where post_id = ?`, postID)
	return comments, err
}

func (c comment) GetByID(ID int) (domain.Comment, error) {
	var comment domain.Comment
	err := c.db.Get(&comment, `select * from comments where id = ?`, ID)
	return comment, err
}
