package store

import (
	"blog/internal/core/port/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	User     repository.UserRepository
	Post     repository.PostRepositoty
	Relation repository.RelationRepository
}

func NewStore(uri string) (Store, error) {
	db, err := sqlx.Connect("mysql", uri)

	store := Store{
		user{db},
		post{db},
		relation{db},
	}
	return store, err
}
