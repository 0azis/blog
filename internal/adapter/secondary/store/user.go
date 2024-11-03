package store

import (
	"blog/internal/core/domain"
	"blog/internal/core/utils"

	"github.com/jmoiron/sqlx"
)

type user struct {
	db *sqlx.DB
}

func (u user) Create(user domain.User) (int, error) {
	sqlResult, err := u.db.Exec(`insert into users (email, username, password, name) values (?, ?, ?, ?)`, user.Email, user.Username, user.Password, user.Name)
	if err != nil {
		return 0, err
	}
	userID, err := sqlResult.LastInsertId()
	return int(userID), err
}

func (u user) GetByID(ID int) (domain.User, error) {
	user := domain.User{}
	rows, err := u.db.Query(`select users.id, users.email, users.username, users.name, users.avatar, users.description, count(distinct posts.id), count(distinct followers.follower_id), count(distinct subscribers.subscriber_id) from users left join posts on users.id = posts.user_id and posts.public = 1 left join followers on users.id = followers.user_id left join subscribers on users.id = subscribers.user_id where users.id = ? group by users.id`, ID)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.Username, &user.Name, &user.Avatar, &user.Description, &user.Counter.Posts, &user.Counter.Followers, &user.Counter.Subscribers)
		if err != nil {
			return user, err
		}
	}
	return user, err
}

func (u user) GetByUsername(username string) (domain.User, error) {
	user := domain.User{}
	rows, err := u.db.Query(`select users.id, users.email, users.username, users.name, users.avatar, users.description, count(distinct posts.id), count(distinct followers.follower_id), count(distinct subscribers.subscriber_id) from users left join posts on users.id = posts.user_id and posts.public = 1 left join followers on users.id = followers.user_id left join subscribers on users.id = subscribers.user_id where users.username = ? group by users.id`, username)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.Username, &user.Name, &user.Avatar, &user.Description, &user.Counter.Posts, &user.Counter.Followers, &user.Counter.Subscribers)
		if err != nil {
			return user, err
		}
	}
	return user, err
}

func (u user) CheckCredentials(email, username string) (domain.User, error) {
	checkedUser := domain.User{}
	err := u.db.Get(&checkedUser, `select * from users where username = ? or email = ?`, username, email)
	return checkedUser, err
}

func (u user) Search(queryMap *utils.QueryMap) ([]*domain.UserCard, error) {
	users := []*domain.UserCard{}
	err := u.db.Select(&users, `select users.id, users.username, users.name, users.avatar from users where lower(username) LIKE lower(?) or lower(name) like lower(?) limit ? offset ?`, queryMap.Queries["q"], queryMap.Queries["q"], queryMap.Pq.Limit, queryMap.Pq.Offset)
	return users, err
}

func (u user) Update(userID int, updatedData domain.UserPatch) (int, error) {
	sqlResult, err := u.db.Exec(`update users set name = ?, avatar = ?, description = ? where id = ?`, updatedData.Name, updatedData.Avatar, updatedData.Description, userID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := sqlResult.RowsAffected()
	return int(rowsAffected), err
}
