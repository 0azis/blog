package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type post struct {
	db *sqlx.DB
}

func (p post) Create(post domain.PostCredentials) (int, error) {
	sqlResult, err := p.db.Exec(`insert into posts (user_id) values (?)`, post.UserID)
	if err != nil {
		return 0, err
	}
	postID, err := sqlResult.LastInsertId()
	return int(postID), err
}

func (p post) Update(postID int, post domain.PostCredentials) (int, error) {
	sqlResult, err := p.db.Exec(`update posts set title = ?, preview = ?, content = ? where user_id = ? and id = ?`, post.Title, post.Preview, post.Content, post.UserID, postID)
	if err != nil {
		return 0, err
	}
	lastID, err := sqlResult.RowsAffected()
	return int(lastID), err
}

func (p post) GetPostsByUser(userID int) ([]domain.UserPost, error) {
	var posts []domain.UserPost
	err := p.db.Select(&posts, `select posts.id, title, date, preview, username, name, avatar, content, count(views.user_id) as views from posts inner join users on posts.user_id = users.id left join views on posts.id = views.post_id where public = 1 and posts.user_id = ? group by posts.id`, userID)
	if err != nil {
		return posts, err
	}

	return posts, err
}

func (p post) GetPostByID(postID int) (domain.UserPost, error) {
	var post domain.UserPost
	err := p.db.Get(&post, `select posts.id, title, date, preview, username, name, avatar, content, count(views.user_id) as views from posts inner join users on posts.user_id = users.id left join views on posts.id = views.post_id where posts.id = ? and public = 1 group by posts.id`, postID)
	// rows, err := p.db.Query(`select posts.id, title, date, preview, username, name, avatar, content from posts inner join users on posts.user_id = users.id where posts.id = ? and public = 1`, postID)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (p post) GetDrafts(userID int) ([]domain.UserPost, error) {
	var drafts []domain.UserPost
	err := p.db.Select(&drafts, `select posts.id, title, date, preview, username, name, avatar, content from posts inner join users on posts.user_id = users.id where user_id = ? and public = 0`, userID)
	return drafts, err
}

func (p post) GetDraft(userID, postID int) (domain.UserPost, error) {
	var draft domain.UserPost
	err := p.db.Get(&draft, `select posts.id, title, date, preview, username, name, avatar, content from posts inner join users on posts.user_id = users.id where posts.id = ? and posts.user_id = ? and public = 0`, postID, userID)
	return draft, err
}

func (p post) Publish(postID, userID int) (int, error) {
	sqlResult, err := p.db.Exec(`update posts set public = 1 where id = ? and user_id = ?`, postID, userID)
	if err != nil {
		return 0, err
	}
	lastID, err := sqlResult.RowsAffected()
	return int(lastID), err
}
