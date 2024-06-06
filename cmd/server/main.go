package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/controllers"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/service"
	"log"
)

func main() {
	r := gin.New()

	appConfig := config.LoadConfig()
	mysqlClient := mysql.NewMysqlClient(appConfig.Database)
	userService := service.NewUserService(mysqlClient)
	handler := controllers.NewUserHandler(userService)
	handler.RegisterRoutes(r)

	if err := r.Run(fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
