package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"database/sql"
	"errors"
	"log/slog"
	nethttp "net/http"

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
		c.JSON(400, utils.JSON{})
		return
	}

	dbUser, err := uc.store.User.CheckCredentials(credentials.Login, credentials.Login)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	if err := utils.Decode([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, utils.JSON{})
		return
	}

	jwts, err := utils.NewJWT(dbUser.ID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", true, true)

	c.JSON(200, utils.JSON{"access_token": jwts.Access})
}

func (uc userControllers) SignUp(c *gin.Context) {
	credentials := domain.SignUpCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	dbUser, err := uc.store.User.CheckCredentials(credentials.Email, credentials.Username)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	validation := domain.ValidateUser(credentials, dbUser)
	if !validation.IsEmail || !validation.IsUsername || !validation.IsPassword {
		c.JSON(409, validation)
		return
	}

	hash, err := utils.Encode([]byte(credentials.Password))
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
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
		c.JSON(500, utils.JSON{})
		return
	}

	jwts, err := utils.NewJWT(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}
	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", true, true)

	c.JSON(201, utils.JSON{"access_token": jwts.Access})
}

func (uc userControllers) Profile(c *gin.Context) {
	userID := utils.ExtractID(c)

	user, err := uc.store.User.GetByID(userID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	user.SetOwnership(userID)
	c.JSON(200, user)
}

func (uc userControllers) GetByUsername(c *gin.Context) {
	userID := utils.ExtractID(c)

	username := c.Param("username")
	account, err := uc.store.User.GetByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	id, err := uc.store.Relation.IsSubcribed(userID, account.ID)
	if id != 0 {
		account.IsSubscribed = true
	}

	if err != nil {
		c.JSON(500, utils.JSON{})
		return
	}

	account.SetOwnership(userID)
	c.JSON(200, account)
}

func (uc userControllers) Search(c *gin.Context) {
	q := map[string]string{}
	err := c.BindQuery(&q)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(400, utils.JSON{})
		return
	}

	queryMap := utils.NewQueryMap(q)
	err = queryMap.SetPaginate()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(400, utils.JSON{})
		return
	}

	search_query := "%" + queryMap.Queries["q"] + "%"
	queryMap.Set("q", search_query)

	queryUsers, err := uc.store.User.Search(queryMap)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, queryUsers)
}

func (uc userControllers) UpdateAccount(c *gin.Context) {
	userID := utils.ExtractID(c)
	var patchCredentials domain.UserPatch
	err := c.ShouldBind(&patchCredentials)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	if !patchCredentials.Validate() {
		c.JSON(400, utils.JSON{})
		return
	}

	rowsAffected, err := uc.store.User.Update(userID, patchCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}
	if rowsAffected == 0 {
		c.JSON(404, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{})
}

func (uc userControllers) Logout(c *gin.Context) {
	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", "", -1, "/", "localhost", true, true)
}

func (uc userControllers) RefreshTokens(c *gin.Context) {
	userID := utils.ExtractID(c)
	jwts, err := utils.NewJWT(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.SetSameSite(nethttp.SameSiteNoneMode)
	c.SetCookie("auth", jwts.Refresh, int(utils.TOKEN_TIME_REFRESH), "/", "localhost", true, true)

	c.JSON(200, utils.JSON{"access_token": jwts.Access})
}

func NewUserControllers(store store.Store) service.UserControllers {
	return userControllers{store}
}
