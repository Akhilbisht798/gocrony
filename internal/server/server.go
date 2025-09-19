package server

import (
	"github.com/akhilbisht798/gocrony/internal/api"
	"github.com/akhilbisht798/gocrony/internal/middleware"
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
	public := s.Router.Group("/api/v1")
	{
		public.POST("/signup", api.EmailPasswordAuthSignUp)
		public.POST("/login", api.EmailPasswordAuthSignIn)
		public.GET("/auth/:provider", api.GetAuthProvider)
		public.GET("/auth/:provider/callback", api.GetAuthCallbackFunction)
		public.GET("/logout/:provider", api.Logout)
	}

	auth := s.Router.Group("/api/v1")
	auth.Use(middleware.AuthMiddleWare())
	{
		auth.POST("/jobs", api.CreateNewJob)
		auth.GET("/jobs", api.GetAllJobs)
		auth.GET("/jobs/:id", api.GetJob)
		auth.PATCH("/jobs/:id", api.UpdateJob)
		auth.DELETE("/jobs/:id", api.DeleteJob)

		auth.POST("/jobs/:id/run", api.RunJob)
		auth.GET("/jobs/:id/logs", api.GetLogs)
	}

	s.Router.Run(addr)
}
