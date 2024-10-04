package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"message"`
	Data   any    `json:"data"`
}

func Error(status int, data any) Response {
	return Response{
		Status: status,
		Msg:    http.StatusText(status),
		Data:   data,
	}
}

func ExtractToken(context *gin.Context) string {
	bearer := context.Request.Header.Get("Authorization")
	jwtToken := strings.Split(bearer, " ")
	return jwtToken[1]
}
