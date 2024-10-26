package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/port/service"
	"blog/internal/core/utils/http"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type relationControllers struct {
	store store.Store
}

func (rc relationControllers) Subscribe(c *gin.Context) {
	userID := http.ExtractID(c)
	value := c.Param("id")
	authorID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}

	err = rc.store.Relation.Subscribe(userID, authorID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(200, http.JSON{})
}

func (rc relationControllers) Followers(c *gin.Context) {
	userID := http.ExtractID(c)
	followers, err := rc.store.Relation.FollowersCount(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(200, http.JSON{"followersCount": followers})
}

func (rc relationControllers) Subscribers(c *gin.Context) {
	userID := http.ExtractID(c)
	subscribers, err := rc.store.Relation.SubscribersCount(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(200, http.JSON{"subscribersCount": subscribers})
}

func NewRelationControllers(store store.Store) service.RelationControllers {
	return relationControllers{store}
}
