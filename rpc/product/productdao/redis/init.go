package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 没有设置密码
		DB:       0,
	})
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Printf("Redis 连接成功，PONG响应: %s", pong)
}
