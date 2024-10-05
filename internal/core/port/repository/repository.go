package repository

import "blog/internal/core/domain"

type UserRepository interface {
	// Create user in database
	Create(user domain.User) (int, error)
	// GetByID get user by its ID
	GetByID(iD int) (domain.User, error)
	// GetByusername get user by its username
	GetByLogin(login string) (domain.User, error)
	// GetByEmail returns userID by its email
	CheckCredentials(email, username string) (domain.User, error)
	// GetByusername get user by its username
	Search(q string, limit, page int) ([]domain.User, error)
}
