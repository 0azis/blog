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

	if userID == authorID {
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
	jwtUserID := utils.ExtractID(c)
	value := c.Param("id")
	userID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	followers, err := rc.store.Relation.Followers(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	for _, follower := range followers {
		follower.SetOwnership(jwtUserID)
	}

	c.JSON(200, utils.JSON{"followers": followers, "followersCount": len(followers)})
}

func (rc relationControllers) Subscribers(c *gin.Context) {
	jwtUserID := utils.ExtractID(c)
	value := c.Param("id")
	userID, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(400, utils.JSON{})
		return
	}

	subscribers, err := rc.store.Relation.Subscribers(userID)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	for _, subscriber := range subscribers {
		subscriber.SetOwnership(jwtUserID)
	}

	c.JSON(200, utils.JSON{"subscribers": subscribers, "subscribersCount": len(subscribers)})
}

func NewRelationControllers(store store.Store) service.RelationControllers {
	return relationControllers{store}
}
