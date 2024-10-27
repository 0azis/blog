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
	Tag      repository.TagRepository
	Comment  repository.CommentRepository
	View     repository.ViewRepository
}

func NewStore(uri string) (Store, error) {
	db, err := sqlx.Connect("mysql", uri)

	store := Store{
		user{db},
		post{db},
		relation{db},
		tag{db},
		comment{db},
		view{db},
	}
	return store, err
}
