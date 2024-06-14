package controllers

import (
	"github.com/gin-gonic/gin"
)

type PictureService struct {
}

func NewPictureService() *PictureService {
	return &PictureService{}
}

func (s *PictureService) RegisterRoutes(router *gin.Engine) {
	router.StaticFile("/avatar", "./public/pictures/avatar")
}
