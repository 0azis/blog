package service

import "github.com/gin-gonic/gin"

type UserControllers interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	Profile(c *gin.Context)
	Search(c *gin.Context)
	RefreshTokens(c *gin.Context)
}

type PostControllers interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetOne(c *gin.Context)
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
