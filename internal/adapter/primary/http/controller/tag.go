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

type tagControllers struct {
	store store.Store
}

func (tc tagControllers) Create(c *gin.Context) {
	var tagCredentials domain.Tag
	err := c.ShouldBind(&tagCredentials)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	rowsAffected, err := tc.store.Tag.Create(tagCredentials)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	if rowsAffected == 0 {
		c.JSON(404, http.Err(404))
		return
	}
	c.JSON(201, http.JSON{})
}

func (tc tagControllers) GetByPostID(c *gin.Context) {
	value := c.Param("id")
	postID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	tags, err := tc.store.Tag.GetByPostID(postID)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(404, http.Err(404))
		return
	}
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(200, tags)
}

func (tc tagControllers) GetByPopularity(c *gin.Context) {
	tags, err := tc.store.Tag.GetByPopularity()
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}
	c.JSON(200, tags)
}

func NewTagControllers(store store.Store) service.TagControllers {
	return &tagControllers{store}
}
