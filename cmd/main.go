package main

import (
	"net/http"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	port := ":" + config.GetEnv("PORT", "8080")
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.POST("/jobs")
	router.GET("/jobs")
	router.GET("/jobs/:id")
	router.PUT("/jobs/:id")
	router.DELETE("/jobs/:id")
	router.GET("/jobs/:id/logs")
	router.POST("/jobs/:id/run")

	router.Run(port)
}
