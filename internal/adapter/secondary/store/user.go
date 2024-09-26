package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type user struct {
	db *sqlx.DB
}

func (u user) Create(user domain.User) (int, error) {
	var userID int
	err := u.db.Get(&userID, `insert into users (firstName, lastName, nick, password) values ($1, $2, $3, $4)`, user.FirstName, user.LastName, user.Nick, user.Password)
	return userID, err
}

func (u user) GetByID(ID int) (domain.User, error) {
	resultUser := domain.User{}
	err := u.db.Get(&resultUser, `select * from users where id = $1`, ID)
	return resultUser, err
}

func (u user) GetByNick(nick string) (domain.User, error) {
	resultUser := domain.User{}
	err := u.db.Get(&resultUser, `select * from users where nick = $1`, nick)
	return resultUser, err
}
