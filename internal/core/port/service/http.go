package service

import "github.com/gin-gonic/gin"

type UserControllers interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	GetByUsername(c *gin.Context)
	Profile(c *gin.Context)
	Search(c *gin.Context)
	UpdateAccount(c *gin.Context)
	Logout(c *gin.Context)
	RefreshTokens(c *gin.Context)
	SignShitIn(c *gin.Context)
}

type PostControllers interface {
	Create(c *gin.Context)
	GetPostsByUser(c *gin.Context)
	MyPosts(c *gin.Context)
	GetByID(c *gin.Context)
	GetDrafts(c *gin.Context)
	GetDraft(c *gin.Context)
	GetPosts(c *gin.Context)
	Publish(c *gin.Context)
	UpdatePost(c *gin.Context)
}

type RelationControllers interface {
	Subscribe(c *gin.Context)
	Followers(c *gin.Context)
	Subscribers(c *gin.Context)
}

type TagControllers interface {
	Create(c *gin.Context)
	GetByPostID(c *gin.Context)
	GetByPopularity(c *gin.Context)
}

type ImageControllers interface {
	Upload(c *gin.Context)
}

type CommentControllres interface {
	NewComment(c *gin.Context)
	GetCommentsByPost(c *gin.Context)
	GetComment(c *gin.Context)
}

type ViewControllers interface {
	ViewsCount(c *gin.Context)
}
