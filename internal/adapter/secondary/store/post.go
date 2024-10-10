package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type post struct {
	db *sqlx.DB
}

func (p post) Create(post domain.PostCredentials) (int, error) {
	sqlResult, err := p.db.Exec(`insert into posts (user_id, content) values (?, ?)`, post.UserID, post.Content)
	if err != nil {
		return 0, err
	}
	postID, err := sqlResult.LastInsertId()
	return int(postID), err
}

func (p post) Update(postID int, post domain.PostCredentials) error {
	_, err := p.db.Query(`update posts set content = ? where user_id = ? and id = ?`, post.Content, post.UserID, postID)
	return err
}

func (p post) GetAll() ([]domain.UserPost, error) {
	var posts []domain.UserPost
	err := p.db.Select(&posts, `select posts.id, username, name, date, content from posts inner join users on posts.user_id = users.id where public = 1`)
	return posts, err
}

func (p post) GetOne(postID int) (domain.UserPost, error) {
	var post domain.UserPost
	err := p.db.Get(&post, `select posts.id, username, name, date, content from posts inner join users on posts.user_id = users.id where posts.id = ?`, postID)
	return post, err
}

func (p post) Publish(postID, userID int) error {
	_, err := p.db.Query("update posts set public = 1 where id = ? and user_id= ?", postID, userID)
	return err
}
