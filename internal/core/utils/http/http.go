package http

import (
	"blog/internal/core/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JSON map[string]any

func Err(statusCode int) JSON {
	return JSON{
		"message": http.StatusText(statusCode),
	}
}

func ErrWithInfo(statusCode int, info any) JSON {
	return JSON{
		"message": http.StatusText(statusCode),
		"info":    info,
	}
}

func ExtractID(context *gin.Context) int {
	bearer := context.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")[1]
	payload, _ := utils.GetIdentity(token)
	return payload.UserID
}
