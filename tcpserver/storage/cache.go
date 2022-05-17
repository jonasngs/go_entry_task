package storage

import (
	"context"
	"log"
	"github.com/go-redis/redis/v8"
	"github.com/jonasngs/go_entry_task/tcpserver/config"
)

type RedisCache struct {
	RDB *redis.Client
}

func InitializeCache() RedisCache {
	redisConfig := config.GetRedisServer()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if _,err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("Error %s when initializing cache\n", err)
	}

	return RedisCache{RDB: rdb}
}
