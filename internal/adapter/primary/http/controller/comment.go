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

type commentControllers struct {
	store store.Store
}

func (cc commentControllers) NewComment(c *gin.Context) {
	var comment domain.Comment
	err := c.ShouldBind(&comment)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}
	userID := utils.ExtractID(c)
	comment.UserID = userID

	err = cc.store.Comment.Create(comment)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(201, utils.JSON{})
}

func (cc commentControllers) GetCommentsByPost(c *gin.Context) {
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(400, utils.JSON{})
		return
	}

	comments, err := cc.store.Comment.GetByPostID(postID)
	if err != nil {
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, comments)
}

func (cc commentControllers) GetComment(c *gin.Context) {
	value := c.Param("id")
	commentID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	comment, err := cc.store.Comment.GetByID(commentID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, utils.JSON{})
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, comment)
}

func NewCommentControllesr(store store.Store) service.CommentControllres {
	return commentControllers{store}
}
