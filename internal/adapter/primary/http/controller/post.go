package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils/http"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type postControllers struct {
	store store.Store
}

func (pc postControllers) Create(c *gin.Context) {
	userID := http.ExtractID(c)
	postCredentials := domain.PostCredentials{}
	err := c.ShouldBind(&postCredentials)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}
	postCredentials.UserID = userID
	postID, err := pc.store.Post.Create(postCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(201, http.JSON{"postID": postID})
}

func (pc postControllers) UpdatePost(c *gin.Context) {
	userID := http.ExtractID(c)
	postCredentials := domain.PostCredentials{}
	err := c.ShouldBind(&postCredentials)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	postCredentials.UserID = userID
	updatedID, err := pc.store.Post.Update(postID, postCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	if updatedID == 0 {
		c.JSON(404, http.Err(404))
		return
	}

	c.JSON(200, http.JSON{})
}

func (pc postControllers) GetPosts(c *gin.Context) {
	posts, err := pc.store.Post.GetPosts()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(200, posts)

}

func (pc postControllers) GetByID(c *gin.Context) {
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	post, err := pc.store.Post.GetPost(postID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, http.Err(404))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(200, post)
}

func (pc postControllers) Publish(c *gin.Context) {
	userID := http.ExtractID(c)
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	post, err := pc.store.Post.GetDraft(postID)
	// fmt.Println("test", *post.Title, *post.Content)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	if !post.Validate() {
		c.JSON(400, http.Err(400))
		return
	}

	publishedID, err := pc.store.Post.Publish(postID, userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	if publishedID == 0 {
		c.JSON(404, http.Err(404))
		return
	}

	c.JSON(200, http.JSON{})
}

func NewPostControllers(store store.Store) service.PostControllers {
	return postControllers{store}
}
