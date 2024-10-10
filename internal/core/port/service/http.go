package service

import "github.com/gin-gonic/gin"

type UserControllers interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	Profile(c *gin.Context)
	Search(c *gin.Context)
}

type PostControllers interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetOne(c *gin.Context)
	Publish(c *gin.Context)
	UpdatePost(c *gin.Context)
}
