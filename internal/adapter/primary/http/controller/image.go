package controller

import (
	"blog/internal/core/port/service"
	"blog/internal/core/utils"
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
		c.JSON(400, utils.JSON{})
		return
	}
	uuid := utils.GenerateUUID()
	filename := uuid + path.Ext(file.Filename)

	err = c.SaveUploadedFile(file, ic.path+"/"+filename)
	if err != nil {
		slog.Error(err.Error())
		c.JSON(500, utils.JSON{})
		return
	}

	c.JSON(201, utils.JSON{"filename": filename})
}

func NewImageControllers(savePath string) service.ImageControllers {
	return imageControllers{savePath}
}
