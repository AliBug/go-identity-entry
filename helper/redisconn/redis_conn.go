package redisconn

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

// NewConn - 初始化 redis 连接
func NewConn(url string) (*redis.Client, error) {
	redisOption, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("Redis url parse error")
	}
	redisClient := redis.NewClient(redisOption)
	return redisClient, nil
}
