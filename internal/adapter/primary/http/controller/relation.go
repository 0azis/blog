package controller

import (
	"blog/internal/adapter/secondary/store"
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type relationControllers struct {
	store store.Store
}

func (rc relationControllers) Subscribe(c *gin.Context) {
	userID := utils.ExtractID(c)
	value := c.Param("id")
	authorID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	err = rc.store.Relation.Subscribe(userID, authorID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{})
}

func (rc relationControllers) Followers(c *gin.Context) {
	userID := utils.ExtractID(c)
	followers, err := rc.store.Relation.FollowersCount(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{"followersCount": followers})
}

func (rc relationControllers) Subscribers(c *gin.Context) {
	userID := utils.ExtractID(c)
	subscribers, err := rc.store.Relation.SubscribersCount(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(200, utils.JSON{"subscribersCount": subscribers})
}

func NewRelationControllers(store store.Store) service.RelationControllers {
	return relationControllers{store}
}
