package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type view struct {
	db *sqlx.DB
}

func (v view) AddView(view domain.View) error {
	_, err := v.db.Query(`insert into views (post_id, user_id) values (?, ?) on duplicate key update post_id=post_id, user_id=user_id`, view.PostID, view.UserID)
	return err
}

func (v view) ViewsCount(postID int) (int, error) {
	var viewCount int
	err := v.db.Get(&viewCount, `select count(user_id) from views where post_id = ?`, postID)
	return viewCount, err
}
