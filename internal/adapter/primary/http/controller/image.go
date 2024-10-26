package controller

import (
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
	"blog/internal/core/utils/http"
	"log/slog"
	"path"

	"github.com/gin-gonic/gin"
)

type imageControllers struct {
	path string
}

func (ic imageControllers) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, http.Err(400))
		return
	}
	uuid := utils.GenerateUUID()
	filename := uuid + path.Ext(file.Filename)

	err = c.SaveUploadedFile(file, ic.path+"/"+filename)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, http.Err(500))
		return
	}

	c.JSON(201, http.JSON{"filename": filename})
}

func NewImageControllers(savePath string) service.ImageControllers {
	return imageControllers{savePath}
}
