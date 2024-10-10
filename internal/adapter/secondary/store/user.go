package store

import (
	"blog/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type user struct {
	db *sqlx.DB
}

func (u user) Create(user domain.User) (int, error) {
	sqlResult, err := u.db.Exec(`insert into users (email, username, password) values (?, ?, ?)`, user.Email, user.Username, user.Password)
	if err != nil {
		return 0, err
	}
	userID, err := sqlResult.LastInsertId()
	return int(userID), err
}

func (u user) GetByID(ID int) (domain.User, error) {
	resultUser := domain.User{}
	err := u.db.Get(&resultUser, `select * from users where id = ?`, ID)
	return resultUser, err
}

func (u user) GetByLogin(login string) (domain.User, error) {
	resultUser := domain.User{}
	err := u.db.Get(&resultUser, `select * from users where username = ? or email = ?`, login, login)
	return resultUser, err
}

func (u user) CheckCredentials(email, username string) (domain.User, error) {
	checkedUser := domain.User{}
	err := u.db.Get(&checkedUser, `select * from users where username = ? or email = ?`, username, email)
	return checkedUser, err
}

func (u user) Search(q string, limit, page int) ([]domain.User, error) {
	resultUser := []domain.User{}
	err := u.db.Select(&resultUser, `select * from users where lower(username) LIKE lower(?) or lower(name) like lower(?) limit ? offset ?`, q, q, limit, page)
	return resultUser, err
}
