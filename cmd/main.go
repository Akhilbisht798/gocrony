package main

import (
	"context"
	"log"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/akhilbisht798/gocrony/internal/auth"
	"github.com/akhilbisht798/gocrony/internal/cache"
	"github.com/akhilbisht798/gocrony/internal/db"
	"github.com/akhilbisht798/gocrony/internal/scheduler"
	"github.com/akhilbisht798/gocrony/internal/server"
	"github.com/akhilbisht798/gocrony/internal/worker"
	"github.com/google/uuid"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	auth.NewAuth()
	err := cache.InitRedisClient()
	if err != nil {
		log.Panic(err)
		return
	}
	go scheduler.Scheduler()
	worker := worker.NewWorker(uuid.NewString())
	go worker.Start(context.Background())

	port := ":" + config.GetEnv("PORT", "8080")

	server := server.NewServer()
	server.Run(port)
}
