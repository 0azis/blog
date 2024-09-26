package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"

	"github.com/gin-gonic/gin"
)

type userControllers struct {
	store *store.Store
}

func (uc userControllers) SignIn(c *gin.Context) {
	credentials := domain.SignInCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	dbUser, err := uc.store.User.GetByNick(credentials.Nick)
	if dbUser.ID == 0 {
		c.JSON(404, gin.H{})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	if err := utils.Decode([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, gin.H{})
		return
	}

	jwt, err := utils.SignJWT(dbUser.ID)
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	c.JSON(200, jwt)
}
func (uc userControllers) SignUp(c *gin.Context) {
	credentials := domain.SignUpCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	dbUserByNick, err := uc.store.User.GetByNick(credentials.Nick)
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	if dbUserByNick.Nick != "" {
		c.JSON(403, gin.H{})
		return
	}

	hash, err := utils.Encode([]byte(credentials.Password))
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	newUser := domain.User{
		FirstName: credentials.FirstName,
		LastName:  credentials.LastName,
		Nick:      credentials.Nick,
		Password:  string(hash),
	}

	userID, err := uc.store.User.Create(newUser)
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	jwt, err := utils.SignJWT(userID)
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	c.JSON(200, jwt)
}

func NewUserControllers(store *store.Store) service.UserControllers {
	return userControllers{store}
}
