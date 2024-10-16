package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/domain"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type tagControllers struct {
	db store.Store
}

func (tc tagControllers) Create(c *gin.Context) {
	var tagCredentials domain.Tag
	err := c.ShouldBind(&tagCredentials)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	rowsAffected, err := tc.db.Tag.Create(tagCredentials)
	fmt.Println(rowsAffected)
	if err != nil {
		c.JSON(500, utils.Error(500, nil))
		return
	}
	if rowsAffected == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	c.JSON(200, utils.Error(200, nil))
}

func (tc tagControllers) GetByPostID(c *gin.Context) {
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.Error(400, nil))
		return
	}

	tags, err := tc.db.Tag.GetByPostID(postID)

	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}

	if tags.PostID == 0 {
		c.JSON(404, utils.Error(404, nil))
		return
	}
	c.JSON(200, utils.Error(200, tags))
}

func (tc tagControllers) GetByPopularity(c *gin.Context) {
	tags, err := tc.db.Tag.GetByPopularity()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.Error(500, nil))
		return
	}
	c.JSON(200, utils.Error(200, tags))
}

func NewTagControllers(store store.Store) service.TagControllers {
	return &tagControllers{store}
}