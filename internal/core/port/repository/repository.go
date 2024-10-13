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

type PostRepositoty interface {
	// Create a post in the database
	Create(post domain.PostCredentials) (int, error)
	// GetAll
	GetAll() ([]domain.UserPost, error)
	// GetOne
	GetOne(postID int) (domain.UserPost, error)
	// Publish
	Publish(postID, userID int) (int, error)
	// Update
	Update(postID int, post domain.PostCredentials) (int, error)
}
