package database

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 初始化Redis客户端
var RC *redis.Client

func ConnectRedis() {
	// 从环境变量读取 Redis 配置
	Redis_Host := os.Getenv("REDIS_HOST")
	Redis_Password := os.Getenv("REDIS_PASSWORD")
	Redis_DB := os.Getenv("REDIS_DB")
	// 把 Redis_DB 转化为 int 类型
	Redis_DBInt, _ := strconv.Atoi(Redis_DB)
	redisClient := redis.NewClient(&redis.Options{
		Addr:         Redis_Host,
		Password:     Redis_Password,
		DB:           Redis_DBInt,
		WriteTimeout: 1 * time.Second, // 写入超时设置为1秒
		ReadTimeout:  1 * time.Second, // 读取超时设置为1秒
		DialTimeout:  3 * time.Second, // 建立连接的超时时间也可以适当调整
		PoolSize:     100,
		MinIdleConns: 3,
		IdleTimeout:  1 * time.Minute,
	})
	RC = redisClient
}
