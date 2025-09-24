package cache

import (
	"context"
	"log"

	"github.com/akhilbisht798/gocrony/config"
	"github.com/redis/go-redis/v9"
)

var Rbd *redis.Client

func InitRedisClient() (error){
	Rbd = redis.NewClient(&redis.Options{
		Addr: config.GetEnv("REDIS_URI", "localhost:6379"),
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB: 0,
		PoolSize: 10,
	})
	pong, err := Rbd.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error connecting to redis: ", err)
		return err
	}
	log.Printf("Connected to redis: %s", pong)
	return nil
}
