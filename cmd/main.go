package main

import (
	"net/http"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/akhilbisht798/gocrony/internal/api"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/scheduler"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db.InitDB()

	port := ":" + config.GetEnv("PORT", "8080")
	router := gin.Default()

	go scheduler.Scheduler()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.POST("/api/v1/jobs", api.CreateNewJob)
	router.GET("/api/v1/jobs", api.GetAllJobs)
	router.GET("/api/v1/jobs/:id", api.GetJob)
	router.PUT("/api/v1/jobs/:id", api.UpdateJob)
	router.DELETE("/api/v1/jobs/:id", api.DeleteJob)
	router.GET("/api/v1/jobs/:id/logs", api.GetLogs)
	router.POST("/api/v1/jobs/:id/run", api.RunJob)

	router.Run(port)
}
