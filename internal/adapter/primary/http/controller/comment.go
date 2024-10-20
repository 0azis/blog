package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
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
		c.JSON(400, utils.Error(400, nil))
		return
	}
	userID := utils.ExtractID(c)
	comment.UserID = userID

	err = cc.store.Comment.Create(comment)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(201, utils.Error(201, nil))
}

func (cc commentControllers) GetCommentsByPost(c *gin.Context) {
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(400, utils.Error(400, nil))
		return
	}

	comments, err := cc.store.Comment.GetByPostID(postID)
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, comments))
}

func (cc commentControllers) GetComment(c *gin.Context) {
	value := c.Param("id")
	commentID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	comment, err := cc.store.Comment.GetByID(commentID)
	if comment.ID == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}

	c.JSON(200, utils.Error(200, comment))
}

func NewCommentControllesr(store store.Store) service.CommentControllres {
	return commentControllers{store}
}
