package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/middlewares"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PictureHandler struct {
	logger      *zap.SugaredLogger
	picturePath string
}

func NewPictureHandler(cfg config.PictureConfig, logger *zap.SugaredLogger) *PictureHandler {
	return &PictureHandler{
		logger:      logger,
		picturePath: cfg.Path,
	}
}

// GetAvatar
//
//	@Summary		Get random picture
//	@Description	Get random picture
//	@Tags			picture
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	middlewares.Response	"ok"
//	@Failure		400	{object}	middlewares.Response	"bad request"
//	@Failure		500	{object}	middlewares.Response	"Internal Server Error"
//	@Router			/user/picture/avatar [get]
func (s *PictureHandler) GetAvatar(ctx *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	files, err := os.ReadDir(s.picturePath)
	if err != nil {
		s.logger.Errorw("Failed to read images directory", "error", err)
		ctx.JSON(http.StatusInternalServerError, middlewares.ErrorResponse(http.StatusBadRequest, "Server internal error"))
		return
	}

	if len(files) == 0 {
		s.logger.Errorw("No images found", "directory", s.picturePath)
		ctx.JSON(http.StatusNotFound, middlewares.ErrorResponse(http.StatusBadRequest, "No images found"))
		return
	}

	// 随机选择一张图片
	randomIndex := rand.Intn(len(files))
	randomFile := files[randomIndex]

	// 获取图片文件的路径
	imagePath := filepath.Join(s.picturePath, randomFile.Name())

	// 返回图片
	ctx.Set("file", imagePath)
}
