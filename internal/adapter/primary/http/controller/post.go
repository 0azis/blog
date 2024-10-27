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
		c.JSON(400, utils.JSON{})
		return
	}
	postCredentials.UserID = userID
	postID, err := pc.store.Post.Create(postCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(201, utils.JSON{"postID": postID})
}

func (pc postControllers) UpdatePost(c *gin.Context) {
	userID := utils.ExtractID(c)
	postCredentials := domain.PostCredentials{}
	err := c.ShouldBind(&postCredentials)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	postCredentials.UserID = userID
	updatedID, err := pc.store.Post.Update(postID, postCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}
	if updatedID == 0 {
		c.JSON(404, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{})
}

func (pc postControllers) GetPosts(c *gin.Context) {
	posts, err := pc.store.Post.GetPosts()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, posts)
}

func (pc postControllers) GetByID(c *gin.Context) {
	userID := utils.ExtractID(c)
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	post, err := pc.store.Post.GetPost(postID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	view := domain.View{
		PostID: post.ID,
		UserID: userID,
	}
	err = pc.store.View.AddView(view)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, post)
}

func (pc postControllers) GetDrafts(c *gin.Context) {
	userID := utils.ExtractID(c)
	drafts, err := pc.store.Post.GetDrafts(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{"drafts": drafts, "draftsCount": len(drafts)})
}

func (pc postControllers) GetDraft(c *gin.Context) {
	userID := utils.ExtractID(c)
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
	}
	draft, err := pc.store.Post.GetDraft(userID, postID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, draft)
}

func (pc postControllers) Publish(c *gin.Context) {
	userID := utils.ExtractID(c)
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	post, err := pc.store.Post.GetDraft(userID, postID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}
	if !post.Validate() {
		c.JSON(400, utils.JSON{})
		return
	}

	publishedID, err := pc.store.Post.Publish(postID, userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}
	if publishedID == 0 {
		c.JSON(404, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{})
}

func NewPostControllers(store store.Store) service.PostControllers {
	return postControllers{store}
}
