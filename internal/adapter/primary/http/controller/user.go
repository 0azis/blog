package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	_ "blog/cmd/docs"

	"github.com/gin-gonic/gin"
)

type userControllers struct {
	store store.Store
}

// SignIn godoc
//
//	@Tags		user
//	@Summary	Sign in to the account
//	@Accept		json
//	@Produce	json
//	@Failure	500	{json}	{"status": 500, "message": "Internal Server Error", "data": null}
//	@Failure	400	{json}	{"status": 400, "message": "Bad Request", "data": null}
//	@Failure	404	{json}	{"status": 404, "message": "Not Found", "data": null}
//	@Failure	401	{json}	{"status": 401, "message": "Unauthorized", "data": null}
//	@Success	200	{json}	{"status": 200, "message": "OK", "data": string}
//	@Router		/user/signin [post]
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

	c.SetCookie("auth", jwts.Refresh, 3600, "/", "localhost", false, false)

	c.JSON(200, utils.Error(200, jwts.Access))
}

// SignUp godoc
//
//	@Tags		user
//	@Summary	Create an account
//	@Accept		json
//	@Produce	json
//	@Failure	500	{json}	{"status": 500, "message": "Internal Server Error", "data": null}
//	@Failure	400	{json}	{"status": 400, "message": "Bad Request", "data": null}
//	@Failure	409	{json}	{"status": 409, "message": "Conflict", "data": null}
//	@Success	201	{json}	{"status": 201, "message": "Created", "data": string}
//	@Router		/user/signup [post]
func (uc userControllers) SignUp(c *gin.Context) {
	credentials := domain.SignUpCredentials{}
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	dbUser, err := uc.store.User.CheckCredentials(credentials.Email, credentials.Username)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	if dbUser.Email == credentials.Email {
		c.JSON(409, utils.Error(
			409, "Email already exists",
		))
		return
	}
	if dbUser.Username == credentials.Username {
		c.JSON(409, utils.Error(
			409, "Username already exists",
		))
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

	c.SetCookie("auth", jwts.Refresh, 3600, "/", "localhost", false, false)

	c.JSON(201, utils.Error(201, jwts.Access))
}

// // Profile godoc
// //
// //	@Tags		user
// //	@Summary	Get a profile
// //	@Accept		json
// //	@Produce	json
// //	@Failure	500	{json}	{"status": 500, "message": "Internal Server Error", "data": null}
// //	@Failure	400	{json}	{"status": 400, "message": "Bad Request", "data": null}
// //	@Failure	409	{json}	{"status": 409, "message": "Conflict", "data": null}
// //	@Success	201	{json}	{"status": 201, "message": "Created", "data": string}
// //	@Router		/user/profile [get]
// func (uc userControllers) Profile(c *gin.Context) {
// 	ID := utils.ExtractID(c)

// 	user, _ := uc.store.User.GetByID(ID)
// 	if user.ID == 0 {
// 		c.JSON(404, utils.Error(404, nil))
// 		return
// 	}

// 	c.JSON(200, utils.Error(200, user))
// }

func (uc userControllers) GetByUsername(c *gin.Context) {
	userID := utils.ExtractID(c)

	username := c.Param("username")
	profile, err := uc.store.User.GetByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	if profile.ID == userID {
		profile.Owner = true
	}

	c.JSON(200, utils.Error(200, profile))
}

//	Search godoc
//
// @Tags		user
// @Summary	Search users from query
// @Accept		json
// @Produce	json
// @Failure	500	{json}	{"status": 500, "message": "Internal Server Error", "data": null}
// @Failure	400	{json}	{"status": 400, "message": "Bad Request", "data": null}
// @Failure	409	{json}	{"status": 409, "message": "Conflict", "data": null}
// @Success	201	{json}	{"status": 201, "message": "Created", "data": string}
// @Router		/user/search [get]
func (uc userControllers) Search(c *gin.Context) {
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

	query_users, err := uc.store.User.Search(search_query, limit, page)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, query_users))
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

	c.JSON(200, utils.Error(200, jwts.Access))
}

func NewUserControllers(store store.Store) service.UserControllers {
	return userControllers{store}
}
