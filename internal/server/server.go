package server

import (
	"net/http"

	"github.com/akhilbisht798/gocrony/internal/api"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
}

func NewServer() *Server {
	router := gin.Default()

	s := &Server{
		Router: router,
	}

	return s
}

func (s *Server) Run(addr string) {
	s.Router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	s.Router.POST("/api/v1/jobs", api.CreateNewJob)
	s.Router.GET("/api/v1/jobs", api.GetAllJobs)
	s.Router.GET("/api/v1/jobs/:id", api.GetJob)
	s.Router.PATCH("/api/v1/jobs/:id", api.UpdateJob)
	s.Router.DELETE("/api/v1/jobs/:id", api.DeleteJob)

	s.Router.POST("/api/v1/jobs/:id/run", api.RunJob)

	s.Router.GET("/api/v1/jobs/:id/logs", api.GetLogs)

	s.Router.GET("/api/v1/auth/:provider/callback", api.GetAuthCallbackFunction)
	s.Router.GET("/api/v1/logout/:provider", api.Logout)
	s.Router.GET("/api/v1/auth/:provider", api.GetAuthProvider)

	s.Router.Run(addr)
}
