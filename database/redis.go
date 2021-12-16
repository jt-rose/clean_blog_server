package database

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Redis connection failed: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println(pong + " Redis connected")
	}
	return rdb
}

var RedisClient = initRedis()