package service

import "github.com/gin-gonic/gin"

type UserControllers interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	Profile(c *gin.Context)
	Search(c *gin.Context)
}
