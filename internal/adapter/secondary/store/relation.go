package store

import "github.com/jmoiron/sqlx"

type relation struct {
	db *sqlx.DB
}

func (r relation) SubscribersCount(userID int) (int, error) {
	var subscribers int
	err := r.db.Get(&subscribers, `select count(id_2) from relations where id_1 = ?`, userID)
	return subscribers, err
}

func (r relation) FollowersCount(userID int) (int, error) {
	var followers int
	err := r.db.Get(&followers, `select count(id_2) from relations where id_2 = ?`, userID)
	return followers, err
}

func (r relation) Subscribe(userID, authorID int) error {
	_, err := r.db.Query(`insert into relations values (?, ?)`, userID, authorID)
	return err
}
