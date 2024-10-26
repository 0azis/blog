package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"blog/internal/core/utils/http"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	nethttp "net/http"
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
		c.JSON(400, http.Err(400))
		return
	}

	dbUser, err := uc.store.User.CheckCredentials(credentials.Login, credentials.Login)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, http.Err(404))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	if err := utils.Decode([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, http.Err(401))
		return
	}

	jwts, err := utils.NewJWT(dbUser.ID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", true, true)

	c.JSON(200, http.JSON{"access_token": jwts.Access})
}

func (uc userControllers) SignUp(c *gin.Context) {
	credentials := domain.SignUpCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	dbUser, err := uc.store.User.CheckCredentials(credentials.Email, credentials.Username)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	validation := domain.ValidateUser(credentials, dbUser)
	if !validation.IsEmail || !validation.IsUsername || !validation.IsPassword {
		c.JSON(409, http.ErrWithInfo(409, http.JSON{"validationInfo": validation}))
		return
	}

	hash, err := utils.Encode([]byte(credentials.Password))
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
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
		c.JSON(500, http.Err(500))
		return
	}

	jwts, err := utils.NewJWT(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", true, true)

	c.JSON(201, http.JSON{"access_token": jwts.Access})
}

func (uc userControllers) Profile(c *gin.Context) {
	userID := http.ExtractID(c)

	user, err := uc.store.User.GetByID(userID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, http.Err(404))
		return
	}
	if err != nil {
		c.JSON(500, http.Err(500))
		return
	}

	user.SetOwnership(userID)
	c.JSON(200, user)
}

func (uc userControllers) GetByUsername(c *gin.Context) {
	userID := http.ExtractID(c)

	username := c.Param("username")
	account, err := uc.store.User.GetByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, http.Err(404))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	account.SetOwnership(userID)
	c.JSON(200, account)
}

func (uc userControllers) Search(c *gin.Context) {
	userID := http.ExtractID(c)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 0 {
		c.JSON(400, http.Err(400))
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	page--
	if err != nil || page < 0 {
		c.JSON(400, http.Err(400))
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
		c.JSON(500, http.Err(500))
		return
	}

	for _, user := range queryUsers {
		fmt.Println(user.ID, userID)
		user.SetOwnership(userID)
	}

	c.JSON(200, queryUsers)
}

func (uc userControllers) UpdateAccount(c *gin.Context) {
	userID := http.ExtractID(c)
	var patchCredentials domain.UserPatch
	err := c.ShouldBind(&patchCredentials)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	if !patchCredentials.Validate() {
		c.JSON(400, http.Err(400))
		return
	}

	rowsAffected, err := uc.store.User.Update(userID, patchCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	if rowsAffected == 0 {
		c.JSON(404, http.Err(404))
		return
	}

	c.JSON(200, http.JSON{})
}

func (uc userControllers) Logout(c *gin.Context) {
	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", "", -1, "/", "localhost", true, true)
}

func (uc userControllers) RefreshTokens(c *gin.Context) {
	userID := http.ExtractID(c)
	jwts, err := utils.NewJWT(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", true, true)

	c.JSON(200, http.JSON{"access_token": jwts.Access})
}

func NewUserControllers(store store.Store) service.UserControllers {
	return userControllers{store}
}
