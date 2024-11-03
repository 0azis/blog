package repository

import (
	"blog/internal/core/domain"
	"blog/internal/core/utils"
)

type UserRepository interface {
	// Create user in database
	Create(user domain.User) (int, error)
	// GetByID get user by its ID
	GetByID(iD int) (domain.User, error)
	// GetByUsername get user by its username
	GetByUsername(username string) (domain.User, error)
	// GetByEmail returns userID by its email
	CheckCredentials(email, username string) (domain.User, error)
	// GetByusername get user by its username
	Search(queryMap *utils.QueryMap) ([]*domain.UserCard, error)
	// Update
	Update(userID int, updatedData domain.UserPatch) (int, error)
}

type PostRepositoty interface {
	// Create a post in the database
	Create(post domain.PostCredentials) (int, error)
	// GetAll
	GetPostsByUser(userID int) ([]*domain.Post, error)
	// GetOne
	GetPostByID(postID int) (domain.Post, error)
	// Get Drafts
	GetDrafts(userID int) ([]*domain.Post, error)
	// GetDraft
	GetDraft(userID, postID int) (domain.Post, error)
	// GetPosts
	GetPosts(queryMap *utils.QueryMap) ([]*domain.Post, error)
	// Publish
	Publish(postID, userID int) (int, error)
	// Update
	Update(postID int, post domain.PostCredentials) (int, error)
}

type RelationRepository interface {
	Subscribers(userID int) ([]*domain.UserCard, error)
	Followers(userID int) ([]*domain.UserCard, error)
	Subscribe(userID, authorID int) error
}

type TagRepository interface {
	Create(tag domain.Tag) (int64, error)
	GetByPostID(postID int) (domain.Tags, error)
	GetByPopularity() (domain.Tags, error)
}

type CommentRepository interface {
	Create(comment domain.CommentCredentials) error
	GetByPostID(postID int) ([]domain.Comment, error)
	GetByID(ID int) (domain.Comment, error)
}

type ViewRepository interface {
	AddView(view domain.View) error
	ViewsCount(postID int) (int, error)
}
