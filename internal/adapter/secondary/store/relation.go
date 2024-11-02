package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type relation struct {
	db *sqlx.DB
}

func (r relation) Subscribers(userID int) ([]*domain.UserCard, error) {
	var subscribers []*domain.UserCard
	err := r.db.Select(&subscribers, `select users.id, users.username, users.name, users.avatar from subscribers left join users on subscribers.subscriber_id = users.id where subscribers.user_id = ?`, userID)
	return subscribers, err
}

func (r relation) Followers(userID int) ([]*domain.UserCard, error) {
	var followers []*domain.UserCard
	err := r.db.Select(&followers, `select users.id, users.username, users.name, users.avatar from followers left join users on followers.follower_id = users.id where followers.user_id = ?`, userID)
	return followers, err
}

func (r relation) Subscribe(userID, authorID int) error {
	_, err := r.db.Exec(`insert into followers values (?, ?)`, authorID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`insert into subscribers values (?, ?)`, userID, authorID)
	return err
}
