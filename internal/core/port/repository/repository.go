package repository

import "blog/internal/core/domain"

type UserRepository interface {
	// Create user in database
	Create(user domain.User) (int, error)
	// GetByID get user by its ID
	GetByID(iD int) (domain.User, error)
	// GetByLogin get user by its username
	GetByLogin(login string) (domain.User, error)
	// GetByUsername get user by its username
	GetByUsername(username string) (domain.User, error)
	// GetByEmail returns userID by its email
	CheckCredentials(email, username string) (domain.User, error)
	// GetByusername get user by its username
	Search(q string, limit, page int) ([]*domain.User, error)
}

type PostRepositoty interface {
	// Create a post in the database
	Create(post domain.PostCredentials) (int, error)
	// GetAll
	GetPosts() ([]domain.UserPost, error)
	// GetOne
	GetPost(postID int) (domain.UserPost, error)
	// Get Drafts
	GetDrafts(userID int) ([]domain.UserPost, error)
	// GetDraft
	GetDraft(postID int) (domain.UserPost, error)
	// Publish
	Publish(postID, userID int) (int, error)
	// Update
	Update(postID int, post domain.PostCredentials) (int, error)
}

type RelationRepository interface {
	SubscribersCount(userID int) (int, error)
	FollowersCount(userID int) (int, error)
	Subscribe(userID, profileID int) error
}

type TagRepository interface {
	Create(tag domain.Tag) (int64, error)
	GetByPostID(postID int) (domain.Tag, error)
	GetByPopularity() (domain.Tags, error)
}

type CommentRepository interface {
	Create(comment domain.Comment) error
	GetByPostID(postID int) ([]domain.Comment, error)
	GetByID(ID int) (domain.Comment, error)
}
