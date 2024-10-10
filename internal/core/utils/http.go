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

func ExtractID(context *gin.Context) int {
	bearer := context.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")[1]
	payload, _ := GetIdentity(token)
	return payload.UserID
}
