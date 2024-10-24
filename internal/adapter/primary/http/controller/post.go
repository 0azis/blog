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

	"github.com/gin-gonic/gin"
)

type postControllers struct {
	store store.Store
}

func (pc postControllers) Create(c *gin.Context) {
	userID := utils.ExtractID(c)
	postCredentials := domain.PostCredentials{}
	err := c.ShouldBind(&postCredentials)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}
	postCredentials.UserID = userID
	postID, err := pc.store.Post.Create(postCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, postID))
}

func (pc postControllers) UpdatePost(c *gin.Context) {
	userID := utils.ExtractID(c)
	postCredentials := domain.PostCredentials{}
	err := c.ShouldBind(&postCredentials)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	postCredentials.UserID = userID
	updatedID, err := pc.store.Post.Update(postID, postCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}
	if updatedID == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}

	c.JSON(200, utils.Error(200, nil))
}

func (pc postControllers) GetPosts(c *gin.Context) {
	posts, err := pc.store.Post.GetAll()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, posts))

}

func (pc postControllers) GetByID(c *gin.Context) {
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	post, err := pc.store.Post.GetOne(postID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, post))
}

func (pc postControllers) Publish(c *gin.Context) {
	userID := utils.ExtractID(c)
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	publishedID, err := pc.store.Post.Publish(postID, userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}
	if publishedID == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}

	c.JSON(200, utils.Error(200, nil))
}

func NewPostControllers(store store.Store) service.PostControllers {
	return postControllers{store}
}
