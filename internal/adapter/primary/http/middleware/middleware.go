package middleware

import (
	"blog/internal/core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		token := strings.Split(bearer, " ")
		if len(token) == 1 {
			c.JSON(401, utils.Error(401, nil))
			c.Abort()
		}
		_, err := utils.GetIdentity(token[1])
		if err != nil {
			c.JSON(401, utils.Error(401, nil))
			c.Abort()
		}
		// status := utils.IsValid(payload)
		// if !status {
		// 	c.JSON(401, utils.Error(401, nil))
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
