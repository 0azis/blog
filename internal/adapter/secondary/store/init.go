package store

import (
	"blog/internal/core/port/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	User repository.UserRepository
}

func NewStore(uri string) (*Store, error) {
	db, err := sqlx.Connect("mysql", uri)

	store := &Store{
		user{db},
	}
	return store, err
}
