package main

import (
	"net/http"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/akhilbisht798/gocrony/internal/api"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	port := ":" + config.GetEnv("PORT", "8080")
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.POST("/jobs", api.CreateNewJob)
	router.GET("/jobs", api.GetAllJobs)
	router.GET("/jobs/:id", api.GetJob)
	router.PUT("/jobs/:id", api.UpdateJob)
	router.DELETE("/jobs/:id", api.DeleteJob)
	router.GET("/jobs/:id/logs", api.GetLogs)
	router.POST("/jobs/:id/run", api.RunJob)

	router.Run(port)
}
