package main

import (
	"github.com/akhilbisht798/gocrony/config"
	"github.com/akhilbisht798/gocrony/internal/auth"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/server"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	auth.NewAuth()

	port := ":" + config.GetEnv("PORT", "8080")

	server := server.NewServer()
	server.Run(port)
}
