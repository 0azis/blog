package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type relation struct {
	db *sqlx.DB
}

func (r relation) Subscribers(userID int) ([]domain.User, error) {
	var subscribers []domain.User
	err := r.db.Select(&subscribers, `select users.* from users inner join relations on users.id = relations.id_2 where relations.id_1 = ?`, userID)
	return subscribers, err
}

func (r relation) Followers(userID int) ([]domain.User, error) {
	var followers []domain.User
	err := r.db.Select(&followers, `select users.* from users inner join relations on users.id = relations.id_1 where relations.id_2 = ?`, userID)
	return followers, err
}

func (r relation) Subscribe(userID, authorID int) error {
	_, err := r.db.Exec(`insert into relations values (?, ?)`, userID, authorID)
	return err
}
