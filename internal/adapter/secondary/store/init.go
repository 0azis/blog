package store

import (
	"blog/internal/core/port/repository"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	User repository.UserRepository
}

func NewStore(uri string) (*Store, error) {
	store := new(Store)

	_, err := sqlx.Connect("pgx", uri)
	// if err != nil {
	// 	return store, err
	// }

	return store, err
}
