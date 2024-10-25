package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	_ "blog/cmd/docs"

	"github.com/gin-gonic/gin"
)

type userControllers struct {
	store store.Store
}

func (uc userControllers) SignIn(c *gin.Context) {
	credentials := domain.SignInCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	dbUser, err := uc.store.User.GetByLogin(credentials.Login)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	if err := utils.Decode([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, utils.Error(401, nil))
		return
	}

	jwts, err := utils.NewJWT(dbUser.ID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", false, false)

	c.JSON(200, utils.Error(200, utils.JSON{"access_token": jwts.Access}))
}

func (uc userControllers) SignUp(c *gin.Context) {
	credentials := domain.SignUpCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	dbUser, err := uc.store.User.CheckCredentials(credentials.Email, credentials.Username)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	validation := domain.ValidateUser(credentials, dbUser)
	if !validation.IsEmail || !validation.IsUsername || !validation.IsPassword {
		c.JSON(409, utils.Error(409, utils.JSON{"validationInfo": validation}))
		return
	}

	hash, err := utils.Encode([]byte(credentials.Password))
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	newUser := domain.User{
		Email:    credentials.Email,
		Username: credentials.Username,
		Password: string(hash),
		Name:     &credentials.Username,
	}

	userID, err := uc.store.User.Create(newUser)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	jwts, err := utils.NewJWT(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", false, false)

	c.JSON(201, utils.Error(201, utils.JSON{"access_token": jwts.Access}))
}

func (uc userControllers) Profile(c *gin.Context) {
	userID := utils.ExtractID(c)

	user, err := uc.store.User.GetByID(userID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}

	user.SetOwnership(userID)
	c.JSON(200, utils.Error(200, utils.JSON{"user": user}))
}

func (uc userControllers) GetByUsername(c *gin.Context) {
	userID := utils.ExtractID(c)

	username := c.Param("username")
	account, err := uc.store.User.GetByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	account.SetOwnership(userID)
	c.JSON(200, utils.Error(200, utils.JSON{"account": account}))
}

func (uc userControllers) Search(c *gin.Context) {
	userID := utils.ExtractID(c)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 0 {
		c.JSON(400, utils.Error(400, nil))
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	page--
	if err != nil || page < 0 {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	if limit == 0 {
		limit = 10
	}

	query := c.Query("q")
	search_query := "%" + query + "%"

	queryUsers, err := uc.store.User.Search(search_query, limit, page)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	for _, user := range queryUsers {
		fmt.Println(user.ID, userID)
		user.SetOwnership(userID)
	}

	c.JSON(200, utils.Error(200, utils.JSON{"queryUsers": queryUsers}))
}

func (uc userControllers) Logout(c *gin.Context) {
	c.SetCookie("auth", "", -1, "/", "localhost", false, false)
}

func (uc userControllers) RefreshTokens(c *gin.Context) {
	userID := utils.ExtractID(c)
	jwts, err := utils.NewJWT(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.SetCookie("auth", jwts.Refresh, 3600, "/", "localhost", false, false)

	c.JSON(200, utils.Error(200, utils.JSON{"access_token": jwts.Access}))
}

func NewUserControllers(store store.Store) service.UserControllers {
	return userControllers{store}
}
