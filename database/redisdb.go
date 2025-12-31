package database

import (
	"time"

	"github.com/go-redis/redis"
)

// 初始化Redis客户端
var RC *redis.Client

func ConnectRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         "43.138.164.205:30007",
		Password:     "1qazWSX!",
		DB:           9,
		WriteTimeout: 1 * time.Second, // 写入超时设置为1秒
		ReadTimeout:  1 * time.Second, // 读取超时设置为1秒
		DialTimeout:  3 * time.Second, // 建立连接的超时时间也可以适当调整
		PoolSize:     100,
		MinIdleConns: 3,
		IdleTimeout:  1 * time.Minute,
	})
	RC = redisClient
}
