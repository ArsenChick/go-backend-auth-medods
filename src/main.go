package main

import (
	"github.com/gin-gonic/gin"

	"github.com/ArsenChick/web-service-gin/controller"
	"github.com/ArsenChick/web-service-gin/middleware/tokenchecker"
	dbservice "github.com/ArsenChick/web-service-gin/services/db"
)

func main() {
	router := gin.Default()

	dbService := dbservice.New()
	defer dbService.CloseConnection()

	ctrl := controller.New(dbService)
	router.POST("/new", ctrl.HandleNewTokenRequest)
	protectedRoutes := router.Group("/refresh", tokenchecker.TokenCheckerMiddleware())
	protectedRoutes.GET("", ctrl.HandleRefreshTokenRequest)

	router.Run(":8080")
}
