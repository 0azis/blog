package middleware

import (
	"blog/internal/core/utils"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	bearer := c.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")
	if len(token) == 1 {
		c.JSON(401, utils.JSON{})
		c.Abort()
		return
	}
	payload, err := utils.GetIdentity(token[1])
	if err != nil {
		c.JSON(401, utils.JSON{})
		c.Abort()
		return
	}
	status := utils.IsValid(payload)
	if !status {
		c.JSON(401, utils.JSON{})
		c.Abort()
		return
	}

	c.Next()
}

func RefreshMiddleware(c *gin.Context) {
	bearer := c.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")
	if len(token) == 1 {
		c.JSON(401, utils.JSON{})
		c.Abort()
		return
	}
	accessPayload, err := utils.GetIdentity(token[1])
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		c.Abort()
		return
	}
	cookieToken, err := c.Cookie("auth")
	if err != nil {
		c.JSON(401, utils.JSON{})
		c.Abort()
		return
	}

	refreshPayload, err := utils.GetIdentity(cookieToken)
	if err != nil {
		c.JSON(500, utils.JSON{})
		c.Abort()
		return
	}

	if refreshPayload.UserID != accessPayload.UserID {
		c.JSON(401, utils.JSON{})
		c.Abort()
		return
	}

	c.Next()
}
