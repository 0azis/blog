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
	rows, err := c.db.Query(`select comments.id, comments.post_id, comments.text, users.id, users.username, users.name, users.avatar from comments inner join users on comments.user_id = users.id where comments.post_id = ?`, postID)
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment domain.Comment
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.Text, &comment.Author.ID, &comment.Author.Username, &comment.Author.Name, &comment.Author.Avatar)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (c comment) GetByID(ID int) (domain.Comment, error) {
	var comment domain.Comment
	rows, err := c.db.Query(`select comments.id, comments.post_id, comments.text, users.id, users.username, users.name, users.avatar from comments inner join users on comments.user_id = users.id where comments.id= ?`, ID)
	if err != nil {
		return comment, err
	}

	for rows.Next() {
		err = rows.Scan(&comment.ID, &comment.PostID, &comment.Text, &comment.Author.ID, &comment.Author.Username, &comment.Author.Name, &comment.Author.Avatar)
	}

	return comment, err
}
