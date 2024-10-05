package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"fmt"
	"strconv"

	_ "blog/cmd/docs"

	"github.com/gin-gonic/gin"
)

type userControllers struct {
	store *store.Store
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

	dbUser, _ := uc.store.User.GetByLogin(credentials.Login)
	if dbUser.ID == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}

	if err := utils.Decode([]byte(dbUser.Password), []byte(credentials.Password)); err != nil {
		c.JSON(401, utils.Error(401, nil))
		return
	}

	accessToken, err := utils.SignJWT(dbUser.ID)
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}
	// // refreshToken, err := utils.SignJWT(dbUser.ID)
	// // if err != nil {
	// // 	c.JSON(500, utils.Error(500, nil))
	// // 	return
	// // }

	// // c.SetCookie("refresh_token", refreshToken, 6, "/", "localhost", false, true)
	// c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)

	c.JSON(200, utils.Error(200, accessToken))
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

	dbUser, _ := uc.store.User.CheckCredentials(credentials.Email, credentials.Username)
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
	}

	hash, err := utils.Encode([]byte(credentials.Password))
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		c.JSON(500, utils.Error(500, nil))
		return
	}

	jwt, err := utils.SignJWT(userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(201, utils.Error(201, jwt))
}

// Profile godoc
//
//	@Tags		user
//	@Summary	Get a profile
//	@Accept		json
//	@Produce	json
//	@Failure	500	{json}	{"status": 500, "message": "Internal Server Error", "data": null}
//	@Failure	400	{json}	{"status": 400, "message": "Bad Request", "data": null}
//	@Failure	409	{json}	{"status": 409, "message": "Conflict", "data": null}
//	@Success	201	{json}	{"status": 201, "message": "Created", "data": string}
//	@Router		/user/profile [get]
func (uc userControllers) Profile(c *gin.Context) {
	clientToken := utils.ExtractToken(c)
	ID, err := utils.GetIdentity(clientToken)
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}

	user, _ := uc.store.User.GetByID(ID)
	if user.ID == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}

	c.JSON(200, utils.Error(200, user))
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
	fmt.Println(query_users)
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, query_users))
}

func NewUserControllers(store *store.Store) service.UserControllers {
	return userControllers{store}
}
